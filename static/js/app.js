// Game state
let gameState = {
    gameId: null,
    playerId: null,
    ws: null,
    currentGame: null
};

// Card emojis
const cardEmojis = {
    'samosa': '🥟',
    'biryani': '🍛',
    'chai': '☕',
    'gulab_jamun': '🍮',
    'paneer_tikka': '🧆',
    'dosa': '🥞',
    'raita': '🥗'
};

const cardNames = {
    'samosa': 'Samosa',
    'biryani': 'Biryani',
    'chai': 'Chai',
    'gulab_jamun': 'Gulab Jamun',
    'paneer_tikka': 'Paneer Tikka',
    'dosa': 'Dosa',
    'raita': 'Raita'
};

// Screen management
function showScreen(screenId) {
    document.querySelectorAll('.screen').forEach(screen => {
        screen.classList.remove('active');
    });
    document.getElementById(screenId).classList.add('active');
}

// Error handling
function showError(message) {
    const errorEl = document.getElementById('errorMessage');
    errorEl.textContent = message;
    errorEl.classList.add('show');
    setTimeout(() => {
        errorEl.classList.remove('show');
    }, 5000);
}

// Home screen
document.getElementById('createGameBtn').addEventListener('click', () => {
    showScreen('createGameScreen');
});

document.getElementById('joinGameBtn').addEventListener('click', () => {
    showScreen('joinGameScreen');
});

// Create game
document.getElementById('confirmCreateBtn').addEventListener('click', async () => {
    const playerName = document.getElementById('createPlayerName').value.trim() || 'Player';

    try {
        const response = await fetch('/api/create', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ player_name: playerName })
        });

        if (!response.ok) {
            throw new Error('Failed to create game');
        }

        const data = await response.json();
        gameState.gameId = data.game_id;
        gameState.playerId = data.player_id;

        connectWebSocket();
        showLobby();
    } catch (error) {
        showError(error.message);
    }
});

document.getElementById('cancelCreateBtn').addEventListener('click', () => {
    showScreen('homeScreen');
});

// Join game
document.getElementById('confirmJoinBtn').addEventListener('click', async () => {
    const gameId = document.getElementById('joinGameId').value.trim();
    const playerName = document.getElementById('joinPlayerName').value.trim() || 'Player';

    if (!gameId) {
        showError('Please enter a game code');
        return;
    }

    try {
        const response = await fetch('/api/join', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                game_id: gameId,
                player_name: playerName
            })
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error || 'Failed to join game');
        }

        const data = await response.json();
        gameState.gameId = data.game_id;
        gameState.playerId = data.player_id;

        connectWebSocket();
        showLobby();
    } catch (error) {
        showError(error.message);
    }
});

document.getElementById('cancelJoinBtn').addEventListener('click', () => {
    showScreen('homeScreen');
});

// Copy game code
document.getElementById('copyCodeBtn').addEventListener('click', () => {
    const code = document.getElementById('gameCodeDisplay').textContent;
    navigator.clipboard.writeText(code);
    showError('Game code copied!');
});

// Start game
document.getElementById('startGameBtn').addEventListener('click', () => {
    sendWebSocketMessage('start_game', {});
});

// Leave lobby
document.getElementById('leaveLobbyBtn').addEventListener('click', () => {
    leaveGame();
});

// Leave game
document.getElementById('leaveGameBtn').addEventListener('click', () => {
    leaveGame();
});

// Back to home
document.getElementById('backToHomeBtn').addEventListener('click', () => {
    showScreen('homeScreen');
});

// WebSocket connection
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws?game_id=${gameState.gameId}&player_id=${gameState.playerId}`;

    gameState.ws = new WebSocket(wsUrl);

    gameState.ws.onopen = () => {
        console.log('WebSocket connected');
        sendWebSocketMessage('get_state', {});
    };

    gameState.ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        handleWebSocketMessage(message);
    };

    gameState.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        showError('Connection error');
    };

    gameState.ws.onclose = () => {
        console.log('WebSocket closed');
    };
}

function sendWebSocketMessage(type, data) {
    if (gameState.ws && gameState.ws.readyState === WebSocket.OPEN) {
        gameState.ws.send(JSON.stringify({ type, data }));
    }
}

function handleWebSocketMessage(message) {
    console.log('Received:', message);

    switch (message.type) {
        case 'game_state':
            updateGameState(message.data);
            break;
        case 'player_joined':
            console.log('Player joined:', message.data);
            break;
        case 'error':
            showError(message.data.message);
            break;
    }
}

// Show lobby
function showLobby() {
    document.getElementById('gameCodeDisplay').textContent = gameState.gameId;
    showScreen('lobbyScreen');
}

// Update game state
function updateGameState(data) {
    gameState.currentGame = data;

    // Update based on game state
    if (data.state === 'waiting') {
        updateLobby(data);
    } else if (data.state === 'playing' || data.state === 'scoring') {
        updateGameScreen(data);
    } else if (data.state === 'finished') {
        showFinalScores(data);
    }
}

// Update lobby
function updateLobby(data) {
    const playerList = document.getElementById('playerList');
    playerList.innerHTML = '';

    document.getElementById('playerCount').textContent = data.players.length;

    data.players.forEach(player => {
        const li = document.createElement('li');
        if (player.id === data.host_id) {
            li.classList.add('host');
        }
        li.textContent = player.name;
        playerList.appendChild(li);
    });

    // Show start button if host
    if (gameState.playerId === data.host_id) {
        document.getElementById('hostControls').style.display = 'block';
        const startBtn = document.getElementById('startGameBtn');
        startBtn.disabled = data.players.length < 2;
    } else {
        document.getElementById('hostControls').style.display = 'none';
    }
}

// Update game screen
function updateGameScreen(data) {
    showScreen('gameScreen');

    // Update header
    document.getElementById('currentRound').textContent = data.round;
    document.getElementById('currentTurn').textContent = data.turn;

    // Find current player
    const currentPlayer = data.players.find(p => p.is_me);
    if (!currentPlayer) return;

    // Update state info
    const stateInfo = document.getElementById('gameStateInfo');
    if (currentPlayer.has_selected) {
        stateInfo.textContent = 'Waiting for other players...';
    } else {
        stateInfo.textContent = 'Select a card from your hand';
    }

    // Update your hand
    document.getElementById('handCount').textContent = currentPlayer.hand ? currentPlayer.hand.length : 0;
    renderHand(currentPlayer.hand || [], currentPlayer.has_selected);

    // Update your played cards
    document.getElementById('yourScore').textContent = currentPlayer.score;
    renderPlayedCards(currentPlayer.played_cards || []);

    // Update other players
    renderOtherPlayers(data.players.filter(p => !p.is_me));
}

// Render hand
function renderHand(hand, hasSelected) {
    const handCards = document.getElementById('handCards');
    handCards.innerHTML = '';

    hand.forEach((card, index) => {
        const cardEl = createCardElement(card);
        if (!hasSelected) {
            cardEl.addEventListener('click', () => selectCard(index));
        } else {
            cardEl.classList.add('disabled');
        }
        handCards.appendChild(cardEl);
    });
}

// Render played cards
function renderPlayedCards(cards) {
    const playedCards = document.getElementById('playedCards');
    playedCards.innerHTML = '';

    cards.forEach(card => {
        const cardEl = createCardElement(card);
        cardEl.classList.add('disabled');
        playedCards.appendChild(cardEl);
    });
}

// Render other players
function renderOtherPlayers(players) {
    const otherPlayers = document.getElementById('otherPlayers');
    otherPlayers.innerHTML = '';

    players.forEach(player => {
        const playerBox = document.createElement('div');
        playerBox.className = 'player-box';
        if (player.has_selected) {
            playerBox.classList.add('selected');
        }

        playerBox.innerHTML = `
            <div class="player-name">${player.name}</div>
            <div class="player-score">Score: ${player.score}</div>
            <div class="player-status">
                Hand: ${player.hand_size} cards |
                Played: ${player.played_cards.length} cards
                ${player.has_selected ? ' ✓' : ''}
            </div>
        `;

        otherPlayers.appendChild(playerBox);
    });
}

// Create card element
function createCardElement(card) {
    const cardEl = document.createElement('div');
    cardEl.className = `game-card ${card.type}`;

    const emoji = document.createElement('div');
    emoji.className = 'card-emoji';
    emoji.textContent = cardEmojis[card.type] || '❓';

    const name = document.createElement('div');
    name.className = 'card-name';
    name.textContent = cardNames[card.type] || card.type;

    cardEl.appendChild(emoji);
    cardEl.appendChild(name);

    // Add value for chai cards
    if (card.type === 'chai' && card.value > 0) {
        const value = document.createElement('div');
        value.className = 'card-value';
        value.textContent = `${card.value} icon${card.value > 1 ? 's' : ''}`;
        cardEl.appendChild(value);
    }

    return cardEl;
}

// Select card
function selectCard(index) {
    sendWebSocketMessage('select_card', { card_index: index });
}

// Show final scores
function showFinalScores(data) {
    const finalScores = document.getElementById('finalScores');
    finalScores.innerHTML = '';

    // Sort players by score
    const sortedPlayers = [...data.players].sort((a, b) => b.score - a.score);

    sortedPlayers.forEach((player, index) => {
        const scoreRow = document.createElement('div');
        scoreRow.className = 'score-row';
        if (index === 0) {
            scoreRow.classList.add('winner');
        }

        scoreRow.innerHTML = `
            <span class="rank">#${index + 1}</span>
            <span>${player.name}</span>
            <span>${player.score} points</span>
        `;

        finalScores.appendChild(scoreRow);
    });

    showScreen('scoresScreen');
}

// Leave game
function leaveGame() {
    if (gameState.ws) {
        gameState.ws.close();
    }
    gameState = {
        gameId: null,
        playerId: null,
        ws: null,
        currentGame: null
    };
    showScreen('homeScreen');
}
