# 🍛 Tiffin Go 🍛

A multiplayer card drafting game inspired by Sushi Go, themed with delicious Indian dishes! Support 2-5 players in real-time web gameplay.

## Game Overview

Tiffin Go is a fast-paced card game where players draft cards over three rounds to create the most valuable collection of Indian dishes. Each round, players simultaneously select one card from their hand, then pass the remaining cards to the next player. Score points based on the combinations of dishes you collect!

## Features

- **Multiplayer**: 2-5 players in real-time
- **WebSocket Communication**: Live game updates
- **Beautiful UI**: Colorful card designs with Indian dish emojis
- **Three Rounds**: Strategic gameplay across multiple rounds
- **7 Card Types**: Each with unique scoring mechanics

## Card Types & Scoring

### 🥟 Samosa
**Scoring**: Pair collection
- Every 2 Samosas = 5 points

### 🍛 Biryani
**Scoring**: Set collection (more is better!)
- 1 Biryani = 1 point
- 2 Biryani = 3 points
- 3 Biryani = 6 points
- 4 Biryani = 10 points
- 5+ Biryani = 15 points

### ☕ Chai
**Scoring**: Competition (most/second most)
- Most Chai icons = 6 points
- Second most Chai icons = 3 points
- Chai cards have 1, 2, or 3 icons

### 🍮 Gulab Jamun
**Scoring**: End game only
- Most Gulab Jamun at game end = +6 points
- Least Gulab Jamun at game end = -6 points (3+ players only)

### 🧆 Paneer Tikka
**Scoring**: Triple collection
- Every 3 Paneer Tikka = 10 points

### 🥞 Dosa
**Special**: Multiplier card
- Triples the value of the next card you play

### 🥗 Raita
**Special**: Action card
- Allows you to play 2 cards in a future turn

## How to Play

### Setup
1. Each player receives a hand of cards (varies by player count)
2. Game lasts 3 rounds

### Each Turn
1. **Select**: Choose one card from your hand
2. **Wait**: Wait for all players to select their card
3. **Reveal**: All selected cards are revealed simultaneously
4. **Pass**: Remaining cards are passed clockwise

### Scoring
- Points are calculated at the end of each round (except Gulab Jamun)
- Gulab Jamun is scored only at the end of the game
- Player with the most points after 3 rounds wins!

## Installation & Running

### Prerequisites
- Go 1.21 or higher

### Build & Run

```bash
# Build the server
go build -o tiffin-go ./cmd/server

# Run the server
./tiffin-go
```

The server will start on `http://localhost:8080`

### Development Mode

```bash
# Run without building
go run ./cmd/server/main.go
```

## Playing the Game

1. **Open your browser** to `http://localhost:8080`
2. **Create a game** or **Join a game** with a 6-digit code
3. **Wait in lobby** for 2-5 players to join
4. **Host starts** the game when ready
5. **Draft cards** by clicking on them in your hand
6. **Watch scores** accumulate over 3 rounds
7. **Winner** is announced at the end!

## Multiplayer Setup

To play with friends:
- Share the 6-digit game code with other players
- All players must connect to the same server
- For remote play, expose port 8080 or deploy to a server

## Project Structure

```
tiffin-go/
├── cmd/
│   └── server/          # Main application entry point
│       └── main.go
├── internal/
│   ├── models/          # Data structures
│   │   ├── card.go      # Card types and deck composition
│   │   ├── player.go    # Player state
│   │   └── game.go      # Game room state
│   ├── game/            # Game logic
│   │   ├── deck.go      # Deck creation and shuffling
│   │   ├── game_logic.go # Turn management
│   │   └── scoring.go   # Scoring algorithms
│   └── server/          # HTTP & WebSocket server
│       ├── game_manager.go  # Multi-game management
│       ├── handlers.go      # HTTP API handlers
│       ├── websocket.go     # WebSocket hub
│       └── ws_handler.go    # WebSocket message handling
├── static/
│   ├── index.html       # Main HTML
│   ├── css/
│   │   └── style.css    # Styling
│   └── js/
│       └── app.js       # Frontend game logic
└── go.mod
```

## API Endpoints

### POST /api/create
Create a new game room
```json
{
  "player_name": "Your Name"
}
```

### POST /api/join
Join an existing game
```json
{
  "game_id": "abc123",
  "player_name": "Your Name"
}
```

### WebSocket /ws
Real-time game communication
```
ws://localhost:8080/ws?game_id=abc123&player_id=xyz789
```

## WebSocket Messages

### Client → Server
- `start_game`: Host starts the game
- `select_card`: Player selects a card
- `get_state`: Request current game state

### Server → Client
- `game_state`: Full game state update
- `player_joined`: New player joined lobby
- `error`: Error message

## Technology Stack

- **Backend**: Go 1.21+
- **WebSocket**: gorilla/websocket
- **Frontend**: Vanilla JavaScript
- **Styling**: CSS3 with gradients and animations

## Future Enhancements

- [ ] Add more card types and variations
- [ ] Implement Raita (play 2 cards) functionality
- [ ] Add sound effects and music
- [ ] Persistent game state (database)
- [ ] Player statistics and leaderboards
- [ ] Multiple game rooms with room browser
- [ ] Mobile responsive improvements
- [ ] AI players for practice mode
- [ ] Tournament mode

## License

MIT License

## Credits

Inspired by the card game "Sushi Go" by Gamewright, adapted with Indian cuisine theme.

---

Enjoy your Tiffin Go game! 🍛☕🥟
