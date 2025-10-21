package models

// Player represents a player in the game
type Player struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Hand        []Card   `json:"hand"`         // Current cards in hand
	PlayedCards []Card   `json:"played_cards"` // Cards played this round
	Score       int      `json:"score"`
	RoundScores []int    `json:"round_scores"` // Score per round
	IsReady     bool     `json:"is_ready"`
	HasSelected bool     `json:"has_selected"` // Has selected card this turn
	DosaActive  bool     `json:"dosa_active"`  // Has active Dosa multiplier
}

// NewPlayer creates a new player
func NewPlayer(id, name string) *Player {
	return &Player{
		ID:          id,
		Name:        name,
		Hand:        []Card{},
		PlayedCards: []Card{},
		Score:       0,
		RoundScores: []int{},
		IsReady:     false,
		HasSelected: false,
		DosaActive:  false,
	}
}
