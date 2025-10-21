package game

import (
	"math/rand"
	"time"

	"github.com/aiplaybookin/tiffin-go/internal/models"
)

// CreateDeck creates a full deck of cards
func CreateDeck() []models.Card {
	deck := []models.Card{}

	// Add Samosa cards
	for i := 0; i < models.DeckComposition[models.Samosa]; i++ {
		deck = append(deck, models.Card{Type: models.Samosa, Value: 0})
	}

	// Add Biryani cards
	for i := 0; i < models.DeckComposition[models.Biryani]; i++ {
		deck = append(deck, models.Card{Type: models.Biryani, Value: 0})
	}

	// Add Chai cards (4 each with 1, 2, 3 icons)
	for i := 1; i <= 3; i++ {
		for j := 0; j < 4; j++ {
			deck = append(deck, models.Card{Type: models.Chai, Value: i})
		}
	}

	// Add Gulab Jamun cards
	for i := 0; i < models.DeckComposition[models.GurabJamun]; i++ {
		deck = append(deck, models.Card{Type: models.GurabJamun, Value: 0})
	}

	// Add Paneer Tikka cards
	for i := 0; i < models.DeckComposition[models.PaneerTikka]; i++ {
		deck = append(deck, models.Card{Type: models.PaneerTikka, Value: 0})
	}

	// Add Dosa cards
	for i := 0; i < models.DeckComposition[models.Dosa]; i++ {
		deck = append(deck, models.Card{Type: models.Dosa, Value: 0})
	}

	// Add Raita cards
	for i := 0; i < models.DeckComposition[models.Raita]; i++ {
		deck = append(deck, models.Card{Type: models.Raita, Value: 0})
	}

	return deck
}

// ShuffleDeck shuffles the deck
func ShuffleDeck(deck []models.Card) []models.Card {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]models.Card, len(deck))
	copy(shuffled, deck)

	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}

// DealCards deals cards to players
func DealCards(game *models.Game) {
	deck := ShuffleDeck(CreateDeck())
	cardsPerPlayer := game.CardsPerHand()

	playerIDs := make([]string, 0, len(game.Players))
	for id := range game.Players {
		playerIDs = append(playerIDs, id)
	}

	cardIndex := 0
	for _, playerID := range playerIDs {
		player := game.Players[playerID]
		player.Hand = []models.Card{}
		for i := 0; i < cardsPerPlayer && cardIndex < len(deck); i++ {
			player.Hand = append(player.Hand, deck[cardIndex])
			cardIndex++
		}
	}

	// Store remaining deck
	if cardIndex < len(deck) {
		game.Deck = deck[cardIndex:]
	} else {
		game.Deck = []models.Card{}
	}
}
