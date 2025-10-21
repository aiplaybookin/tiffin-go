package game

import (
	"errors"

	"github.com/aiplaybookin/tiffin-go/internal/models"
)

// StartGame initializes a new game
func StartGame(game *models.Game) error {
	if !game.CanStart() {
		return errors.New("cannot start game: need 2-5 players")
	}

	game.State = models.StatePlaying
	game.Round = 1
	game.Turn = 1

	// Initialize all players
	for _, player := range game.Players {
		player.IsReady = true
		player.HasSelected = false
		player.PlayedCards = []models.Card{}
		player.DosaActive = false
	}

	DealCards(game)
	return nil
}

// SelectCard handles a player selecting a card from their hand
func SelectCard(game *models.Game, playerID string, cardIndex int) error {
	player, exists := game.Players[playerID]
	if !exists {
		return errors.New("player not found")
	}

	if player.HasSelected {
		return errors.New("player has already selected a card this turn")
	}

	if cardIndex < 0 || cardIndex >= len(player.Hand) {
		return errors.New("invalid card index")
	}

	if game.State != models.StatePlaying {
		return errors.New("game is not in playing state")
	}

	// Remove card from hand and add to played cards
	selectedCard := player.Hand[cardIndex]
	player.Hand = append(player.Hand[:cardIndex], player.Hand[cardIndex+1:]...)
	player.PlayedCards = append(player.PlayedCards, selectedCard)
	player.HasSelected = true

	// Check if Dosa was played (will be active for next card)
	if selectedCard.Type == models.Dosa {
		player.DosaActive = true
	}

	return nil
}

// PassHands rotates hands clockwise to the next player
func PassHands(game *models.Game) error {
	if !game.AllPlayersSelected() {
		return errors.New("not all players have selected")
	}

	// Get ordered list of player IDs
	playerIDs := make([]string, 0, len(game.Players))
	for id := range game.Players {
		playerIDs = append(playerIDs, id)
	}

	if len(playerIDs) == 0 {
		return errors.New("no players in game")
	}

	// Save current hands
	hands := make(map[string][]models.Card)
	for _, id := range playerIDs {
		hands[id] = game.Players[id].Hand
	}

	// Pass hands clockwise
	for i := 0; i < len(playerIDs); i++ {
		currentPlayerID := playerIDs[i]
		nextPlayerID := playerIDs[(i+1)%len(playerIDs)]
		game.Players[nextPlayerID].Hand = hands[currentPlayerID]
	}

	// Reset selection flags
	for _, player := range game.Players {
		player.HasSelected = false
	}

	game.Turn++

	// Check if round is over (hands are empty or only 1 card left)
	roundOver := true
	for _, player := range game.Players {
		if len(player.Hand) > 1 {
			roundOver = false
			break
		}
	}

	if roundOver {
		return EndRound(game)
	}

	return nil
}

// EndRound handles end of round scoring
func EndRound(game *models.Game) error {
	game.State = models.StateScoring

	// Calculate scores for this round
	ScoreRound(game)

	// Reset for next round or end game
	if game.Round >= 3 {
		game.State = models.StateFinished
		FinalScoring(game)
	} else {
		game.Round++
		game.Turn = 1

		// Clear played cards except Gulab Jamun (pudding)
		for _, player := range game.Players {
			puddings := []models.Card{}
			for _, card := range player.PlayedCards {
				if card.Type == models.GurabJamun {
					puddings = append(puddings, card)
				}
			}
			player.PlayedCards = puddings
			player.DosaActive = false
			player.HasSelected = false
		}

		// Deal new cards
		DealCards(game)
		game.State = models.StatePlaying
	}

	return nil
}
