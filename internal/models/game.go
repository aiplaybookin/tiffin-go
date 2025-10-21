package models

import (
	"time"
)

// GameState represents the current state of the game
type GameState string

const (
	StateWaiting  GameState = "waiting"   // Waiting for players
	StatePlaying  GameState = "playing"   // Game in progress
	StateScoring  GameState = "scoring"   // Between rounds, showing scores
	StateFinished GameState = "finished"  // Game complete
)

// Game represents a game room
type Game struct {
	ID          string               `json:"id"`
	Players     map[string]*Player   `json:"players"`      // PlayerID -> Player
	State       GameState            `json:"state"`
	Round       int                  `json:"round"`        // 1, 2, or 3
	Turn        int                  `json:"turn"`         // Current turn in round
	Deck        []Card               `json:"-"`            // Remaining cards in deck
	HostID      string               `json:"host_id"`
	CreatedAt   time.Time            `json:"created_at"`
	MaxPlayers  int                  `json:"max_players"`
	MinPlayers  int                  `json:"min_players"`
}

// NewGame creates a new game room
func NewGame(id string, hostID string) *Game {
	return &Game{
		ID:         id,
		Players:    make(map[string]*Player),
		State:      StateWaiting,
		Round:      0,
		Turn:       0,
		Deck:       []Card{},
		HostID:     hostID,
		CreatedAt:  time.Now(),
		MaxPlayers: 5,
		MinPlayers: 2,
	}
}

// CanStart checks if the game can be started
func (g *Game) CanStart() bool {
	playerCount := len(g.Players)
	return playerCount >= g.MinPlayers && playerCount <= g.MaxPlayers && g.State == StateWaiting
}

// AllPlayersSelected checks if all players have selected a card
func (g *Game) AllPlayersSelected() bool {
	for _, player := range g.Players {
		if !player.HasSelected {
			return false
		}
	}
	return true
}

// CardsPerHand returns number of cards dealt based on player count
func (g *Game) CardsPerHand() int {
	playerCount := len(g.Players)
	switch playerCount {
	case 2:
		return 10
	case 3:
		return 9
	case 4:
		return 8
	case 5:
		return 7
	default:
		return 8
	}
}
