package models

// CardType represents the type of card
type CardType string

const (
	Samosa       CardType = "samosa"        // Pairs score 5 points
	Biryani      CardType = "biryani"       // Set collection: 1=1, 2=3, 3=6, 4=10, 5+=15 points
	Chai         CardType = "chai"          // Most chai icons = 6 points, 2nd = 3 points
	GurabJamun   CardType = "gulab_jamun"   // Pudding scoring at end (most=6, least=-6)
	PaneerTikka  CardType = "paneer_tikka"  // 3 cards = 10 points
	Dosa         CardType = "dosa"          // Multiplier: triple next card value
	Raita        CardType = "raita"         // Play 2 cards in future turn
)

// Card represents a single game card
type Card struct {
	Type  CardType `json:"type"`
	Value int      `json:"value"` // For Chai - number of icons (1-3)
}

// CardCount represents how many of each card type in the deck
var DeckComposition = map[CardType]int{
	Samosa:      14,
	Biryani:     14,
	Chai:        12, // 4 cards each with 1, 2, 3 icons
	GurabJamun:  10,
	PaneerTikka: 14,
	Dosa:        6,
	Raita:       4,
}
