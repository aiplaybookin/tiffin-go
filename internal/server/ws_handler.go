package server

import (
	"encoding/json"
	"log"

	"github.com/aiplaybookin/tiffin-go/internal/game"
	"github.com/aiplaybookin/tiffin-go/internal/models"
)

// WSHandler handles WebSocket messages
type WSHandler struct {
	hub         *Hub
	gameManager *GameManager
}

// NewWSHandler creates a new WebSocket handler
func NewWSHandler(hub *Hub, gm *GameManager) *WSHandler {
	return &WSHandler{
		hub:         hub,
		gameManager: gm,
	}
}

// WSMessage represents an incoming WebSocket message
type WSMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// SelectCardData represents data for selecting a card
type SelectCardData struct {
	CardIndex int `json:"card_index"`
}

// HandleMessage processes incoming WebSocket messages
func (wh *WSHandler) HandleMessage(client *Client, message []byte) {
	var msg WSMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	switch msg.Type {
	case "select_card":
		wh.handleSelectCard(client, msg.Data)
	case "start_game":
		wh.handleStartGame(client)
	case "get_state":
		wh.handleGetState(client)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// handleSelectCard processes card selection
func (wh *WSHandler) handleSelectCard(client *Client, data json.RawMessage) {
	var selectData SelectCardData
	if err := json.Unmarshal(data, &selectData); err != nil {
		log.Printf("Error unmarshaling select card data: %v", err)
		return
	}

	g, err := wh.gameManager.GetGame(client.GameID)
	if err != nil {
		log.Printf("Game not found: %v", err)
		return
	}

	// Select card
	err = game.SelectCard(g, client.ID, selectData.CardIndex)
	if err != nil {
		log.Printf("Error selecting card: %v", err)
		wh.sendError(client, err.Error())
		return
	}

	// Broadcast game state
	wh.broadcastGameState(client.GameID)

	// If all players selected, pass hands
	if g.AllPlayersSelected() {
		err = game.PassHands(g)
		if err != nil {
			log.Printf("Error passing hands: %v", err)
			return
		}
		wh.broadcastGameState(client.GameID)
	}
}

// handleStartGame starts the game
func (wh *WSHandler) handleStartGame(client *Client) {
	g, err := wh.gameManager.GetGame(client.GameID)
	if err != nil {
		log.Printf("Game not found: %v", err)
		return
	}

	// Only host can start
	if g.HostID != client.ID {
		wh.sendError(client, "only host can start the game")
		return
	}

	err = game.StartGame(g)
	if err != nil {
		log.Printf("Error starting game: %v", err)
		wh.sendError(client, err.Error())
		return
	}

	wh.broadcastGameState(client.GameID)
}

// handleGetState sends current game state to client
func (wh *WSHandler) handleGetState(client *Client) {
	wh.broadcastGameState(client.GameID)
}

// broadcastGameState sends game state to all clients
func (wh *WSHandler) broadcastGameState(gameID string) {
	g, err := wh.gameManager.GetGame(gameID)
	if err != nil {
		return
	}

	// Create sanitized state for each player (hide other players' hands)
	wh.hub.mu.RLock()
	for _, client := range wh.hub.clients {
		if client.GameID == gameID {
			state := createPlayerGameState(g, client.ID)
			wh.sendToClient(client, "game_state", state)
		}
	}
	wh.hub.mu.RUnlock()
}

// sendToClient sends a message to a specific client
func (wh *WSHandler) sendToClient(client *Client, messageType string, data interface{}) {
	message := map[string]interface{}{
		"type": messageType,
		"data": data,
	}

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	select {
	case client.Send <- jsonMsg:
	default:
		log.Printf("Client send buffer full")
	}
}

// sendError sends an error message to client
func (wh *WSHandler) sendError(client *Client, errorMsg string) {
	wh.sendToClient(client, "error", map[string]string{"message": errorMsg})
}

// createPlayerGameState creates a game state with hidden information for other players
func createPlayerGameState(g *models.Game, playerID string) map[string]interface{} {
	// Create players list with hidden hands
	players := make([]map[string]interface{}, 0)
	for id, p := range g.Players {
		playerData := map[string]interface{}{
			"id":            p.ID,
			"name":          p.Name,
			"score":         p.Score,
			"round_scores":  p.RoundScores,
			"has_selected":  p.HasSelected,
			"played_cards":  p.PlayedCards,
			"hand_size":     len(p.Hand),
		}

		// Only show full hand to the player themselves
		if id == playerID {
			playerData["hand"] = p.Hand
			playerData["is_me"] = true
		} else {
			playerData["is_me"] = false
		}

		players = append(players, playerData)
	}

	return map[string]interface{}{
		"id":      g.ID,
		"state":   g.State,
		"round":   g.Round,
		"turn":    g.Turn,
		"host_id": g.HostID,
		"players": players,
	}
}
