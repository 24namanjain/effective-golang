# Nakama - Super Simple Guide

## What is Nakama?

Nakama is a **game server**. It's like a middleman between players in a multiplayer game.

## How it works (in 3 steps)

1. **Players connect** to Nakama server
2. **Players send requests** (like "create game", "join game", "submit score")
3. **Server responds** and tells other players what happened

## What you can do with it

- **Chat** - Players can message each other
- **Leaderboards** - Track high scores
- **Multiplayer games** - Real-time game sessions
- **User accounts** - Store player data

## Quick Start

### 1. Start the server
```bash
docker compose up -d
```

### 2. Test it works
```bash
node test-client.js
```

### 3. Try the game example
```bash
node examples/client-example.js
```

## Files Explained

- `docker-compose.yml` - Starts Nakama and database
- `test-client.js` - Simple test to check if it works
- `examples/simple-game.go` - Game functions (runs on server)
- `examples/client-example.js` - Game client (runs on player's device)
- `build-module.sh` - Builds the game functions

## Simple Example

**Server (Go)**: Defines what players can do
```go
// When a player calls "create_game"
func createGame(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    return `{"success": true, "game_id": "game123"}`, nil
}
```

**Client (JavaScript)**: Player calls the function
```javascript
// Player wants to create a game
const result = await client.rpc(session, 'create_game', {name: 'My Game'});
console.log(result); // {"success": true, "game_id": "game123"}
```

## That's it!

- ✅ **Server running** - Nakama handles multiplayer
- ✅ **Database ready** - Stores player data
- ✅ **Game functions** - Players can create/join games
- ✅ **Real-time** - Instant communication between players

**Next**: Read the examples and modify them to build your own game!
