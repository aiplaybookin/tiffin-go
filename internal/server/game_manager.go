package server

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/aiplaybookin/tiffin-go/internal/models"
)

// GameManager manages all active games
type GameManager struct {
	games map[string]*models.Game
	mu    sync.RWMutex
}

// NewGameManager creates a new game manager
func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[string]*models.Game),
	}
}

// CreateGame creates a new game room
func (gm *GameManager) CreateGame(hostID, hostName string) (*models.Game, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	gameID := generateGameID()
	game := models.NewGame(gameID, hostID)

	// Add host as first player
	player := models.NewPlayer(hostID, hostName)
	game.Players[hostID] = player

	gm.games[gameID] = game
	return game, nil
}

// JoinGame adds a player to an existing game
func (gm *GameManager) JoinGame(gameID, playerID, playerName string) (*models.Game, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, exists := gm.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	if game.State != models.StateWaiting {
		return nil, errors.New("game already started")
	}

	if len(game.Players) >= game.MaxPlayers {
		return nil, errors.New("game is full")
	}

	// Check if player already in game
	if _, exists := game.Players[playerID]; exists {
		return game, nil // Already joined
	}

	player := models.NewPlayer(playerID, playerName)
	game.Players[playerID] = player

	return game, nil
}

// LeaveGame removes a player from a game
func (gm *GameManager) LeaveGame(gameID, playerID string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, exists := gm.games[gameID]
	if !exists {
		return errors.New("game not found")
	}

	delete(game.Players, playerID)

	// If no players left or host left while waiting, delete game
	if len(game.Players) == 0 || (game.HostID == playerID && game.State == models.StateWaiting) {
		delete(gm.games, gameID)
	}

	return nil
}

// GetGame retrieves a game by ID
func (gm *GameManager) GetGame(gameID string) (*models.Game, error) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	game, exists := gm.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	return game, nil
}

// generateGameID creates a random 6-character game ID
func generateGameID() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
