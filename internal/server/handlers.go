package server

import (
	"encoding/json"
	"log"
	"net/http"

)

// Server represents the HTTP server
type Server struct {
	hub         *Hub
	gameManager *GameManager
	wsHandler   *WSHandler
}

// NewServer creates a new server
func NewServer() *Server {
	hub := NewHub()
	gameManager := NewGameManager()
	wsHandler := NewWSHandler(hub, gameManager)

	return &Server{
		hub:         hub,
		gameManager: gameManager,
		wsHandler:   wsHandler,
	}
}

// Start starts the server
func (s *Server) Start() {
	go s.hub.Run()
}

// CreateGameRequest represents a request to create a game
type CreateGameRequest struct {
	PlayerName string `json:"player_name"`
}

// CreateGameResponse represents the response from creating a game
type CreateGameResponse struct {
	GameID   string `json:"game_id"`
	PlayerID string `json:"player_id"`
}

// JoinGameRequest represents a request to join a game
type JoinGameRequest struct {
	GameID     string `json:"game_id"`
	PlayerName string `json:"player_name"`
}

// JoinGameResponse represents the response from joining a game
type JoinGameResponse struct {
	GameID   string `json:"game_id"`
	PlayerID string `json:"player_id"`
}

// HandleCreateGame handles creating a new game
func (s *Server) HandleCreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.PlayerName == "" {
		req.PlayerName = "Player"
	}

	// Generate player ID
	playerID := generateGameID()

	game, err := s.gameManager.CreateGame(playerID, req.PlayerName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CreateGameResponse{
		GameID:   game.ID,
		PlayerID: playerID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleJoinGame handles joining an existing game
func (s *Server) HandleJoinGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req JoinGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.GameID == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	if req.PlayerName == "" {
		req.PlayerName = "Player"
	}

	// Generate player ID
	playerID := generateGameID()

	game, err := s.gameManager.JoinGame(req.GameID, playerID, req.PlayerName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Broadcast player joined
	s.hub.BroadcastToGame(game.ID, "player_joined", map[string]string{
		"player_id":   playerID,
		"player_name": req.PlayerName,
	})

	resp := JoinGameResponse{
		GameID:   game.ID,
		PlayerID: playerID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleWebSocket handles WebSocket connections
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("game_id")
	playerID := r.URL.Query().Get("player_id")

	if gameID == "" || playerID == "" {
		http.Error(w, "game_id and player_id required", http.StatusBadRequest)
		return
	}

	// Verify game and player exist
	game, err := s.gameManager.GetGame(gameID)
	if err != nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	if _, exists := game.Players[playerID]; !exists {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:     playerID,
		GameID: gameID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	s.hub.register <- client

	// Start client pumps
	go client.writePump()
	go client.readPump(s.hub, s.wsHandler)

	// Send initial game state
	s.wsHandler.broadcastGameState(gameID)
}
