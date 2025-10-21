package game

import (
	"github.com/aiplaybookin/tiffin-go/internal/models"
)

// ScoreRound calculates scores for all players at end of round
func ScoreRound(game *models.Game) {
	for _, player := range game.Players {
		roundScore := 0

		// Count cards by type
		cardCounts := make(map[models.CardType]int)
		chaiIcons := 0

		for _, card := range player.PlayedCards {
			if card.Type != models.GurabJamun { // Don't count pudding in round score
				cardCounts[card.Type]++
				if card.Type == models.Chai {
					chaiIcons += card.Value
				}
			}
		}

		// Score Samosa (pairs = 5 points)
		samosaPairs := cardCounts[models.Samosa] / 2
		roundScore += samosaPairs * 5

		// Score Biryani (set collection: 1=1, 2=3, 3=6, 4=10, 5+=15)
		biryaniCount := cardCounts[models.Biryani]
		switch {
		case biryaniCount == 1:
			roundScore += 1
		case biryaniCount == 2:
			roundScore += 3
		case biryaniCount == 3:
			roundScore += 6
		case biryaniCount == 4:
			roundScore += 10
		case biryaniCount >= 5:
			roundScore += 15
		}

		// Score Paneer Tikka (3 = 10 points)
		paneerSets := cardCounts[models.PaneerTikka] / 3
		roundScore += paneerSets * 10

		// Chai scoring (most icons = 6, second most = 3)
		// Will be calculated after all players counted
		player.RoundScores = append(player.RoundScores, roundScore)
	}

	// Score Chai (most/second most)
	scoreChaiForPlayers(game)

	// Update total scores
	for _, player := range game.Players {
		if len(player.RoundScores) > 0 {
			lastScore := player.RoundScores[len(player.RoundScores)-1]
			player.Score += lastScore
		}
	}
}

// scoreChaiForPlayers awards points for most/second most chai icons
func scoreChaiForPlayers(game *models.Game) {
	type playerChai struct {
		playerID string
		icons    int
	}

	chaiCounts := []playerChai{}
	for id, player := range game.Players {
		icons := 0
		for _, card := range player.PlayedCards {
			if card.Type == models.Chai {
				icons += card.Value
			}
		}
		if icons > 0 {
			chaiCounts = append(chaiCounts, playerChai{id, icons})
		}
	}

	if len(chaiCounts) == 0 {
		return
	}

	// Sort by icons descending
	for i := 0; i < len(chaiCounts); i++ {
		for j := i + 1; j < len(chaiCounts); j++ {
			if chaiCounts[j].icons > chaiCounts[i].icons {
				chaiCounts[i], chaiCounts[j] = chaiCounts[j], chaiCounts[i]
			}
		}
	}

	// Award points for most
	mostIcons := chaiCounts[0].icons
	mostPlayers := []string{}
	for _, pc := range chaiCounts {
		if pc.icons == mostIcons {
			mostPlayers = append(mostPlayers, pc.playerID)
		}
	}

	pointsPerPlayer := 6 / len(mostPlayers)
	for _, playerID := range mostPlayers {
		if len(game.Players[playerID].RoundScores) > 0 {
			game.Players[playerID].RoundScores[len(game.Players[playerID].RoundScores)-1] += pointsPerPlayer
		}
	}

	// Award points for second most (if not tied for first)
	if len(chaiCounts) > len(mostPlayers) {
		secondMostIcons := 0
		for _, pc := range chaiCounts {
			if pc.icons < mostIcons {
				secondMostIcons = pc.icons
				break
			}
		}

		if secondMostIcons > 0 {
			secondPlayers := []string{}
			for _, pc := range chaiCounts {
				if pc.icons == secondMostIcons {
					secondPlayers = append(secondPlayers, pc.playerID)
				}
			}

			pointsPerPlayer := 3 / len(secondPlayers)
			for _, playerID := range secondPlayers {
				if len(game.Players[playerID].RoundScores) > 0 {
					game.Players[playerID].RoundScores[len(game.Players[playerID].RoundScores)-1] += pointsPerPlayer
				}
			}
		}
	}
}

// FinalScoring calculates Gulab Jamun (pudding) points at game end
func FinalScoring(game *models.Game) {
	type playerPudding struct {
		playerID string
		count    int
	}

	puddingCounts := []playerPudding{}
	for id, player := range game.Players {
		count := 0
		for _, card := range player.PlayedCards {
			if card.Type == models.GurabJamun {
				count++
			}
		}
		puddingCounts = append(puddingCounts, playerPudding{id, count})
	}

	if len(puddingCounts) == 0 {
		return
	}

	// Sort by count
	for i := 0; i < len(puddingCounts); i++ {
		for j := i + 1; j < len(puddingCounts); j++ {
			if puddingCounts[j].count > puddingCounts[i].count {
				puddingCounts[i], puddingCounts[j] = puddingCounts[j], puddingCounts[i]
			}
		}
	}

	// Most pudding gets +6
	mostCount := puddingCounts[0].count
	for _, pp := range puddingCounts {
		if pp.count == mostCount {
			game.Players[pp.playerID].Score += 6
		}
	}

	// Least pudding gets -6 (only if 3+ players)
	if len(game.Players) >= 3 {
		leastCount := puddingCounts[len(puddingCounts)-1].count
		if leastCount < mostCount {
			for _, pp := range puddingCounts {
				if pp.count == leastCount {
					game.Players[pp.playerID].Score -= 6
				}
			}
		}
	}
}
