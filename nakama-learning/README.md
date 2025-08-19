# Nakama Learning - Simplified

A simple way to learn Nakama game server.

## What is Nakama?

Nakama is a **game server** that handles multiplayer functionality for games. Think of it as a backend service that manages:
- Player connections
- Real-time messaging
- Leaderboards
- User data storage

## Quick Start (3 Steps)

### 1. Start the Server
```bash
docker compose up -d
```

### 2. Test It Works
```bash
node test-client.js
```

### 3. View Logs
```bash
docker compose logs -f nakama
```

## What's Running?

- **Nakama Server** - Handles game logic (port 7352)
- **Database** - Stores player data (port 26258)

## Simple Example

**Server (Go)**: Defines game functions
```go
// This function runs on the server
func createGame(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    // Create a new game
    return `{"success": true, "game_id": "game123"}`, nil
}
```

**Client (JavaScript)**: Calls the function
```javascript
// This runs on the player's device
const result = await client.rpc(session, 'create_game', {name: 'My Game'});
console.log(result); // {"success": true, "game_id": "game123"}
```

## How It Works

1. **Player connects** to Nakama server
2. **Player calls a function** (like "create game")
3. **Server runs the function** and responds
4. **Other players get notified** in real-time

## Common Use Cases

- **Chat system** - Players can message each other
- **Leaderboards** - Track high scores
- **Multiplayer games** - Real-time game sessions
- **User profiles** - Store player data

## Files Explained

- `docker-compose.yml` - Starts Nakama and database
- `test-client.js` - Tests if everything works
- `examples/simple-game.go` - Shows how to write game functions
- `examples/client-example.js` - Shows how to connect from a game

## Next Steps

1. **Read the examples** to see how it works
2. **Modify the code** to add your own features
3. **Build your game** using the patterns shown

---

**That's it!** Much simpler than before. Nakama is just a tool to help you build multiplayer games.
