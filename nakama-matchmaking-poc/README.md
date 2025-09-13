# Nakama Matchmaking POC

A clean proof of concept demonstrating Nakama's built-in matchmaking capabilities with real-time multiplayer matches.

## 🎯 What This POC Demonstrates

- **Built-in Matchmaker**: Uses Nakama's native matchmaking system
- **Real-time Matches**: Players are matched and join live game sessions
- **Match State Management**: Tracks player readiness, scores, and game state
- **Multi-client Testing**: Simulates multiple players joining and playing

## 🚀 Quick Start

### 1. Start Nakama Server

```bash
# Start Nakama and database
docker-compose up -d

# Check if services are running
docker-compose ps
```

This starts:
- **Nakama server** on ports 7350 (HTTP) and 7351 (gRPC)
- **CockroachDB** on port 26257

### 2. Install Dependencies

```bash
# Install Node.js dependencies
npm install
```

### 3. Build Go Module

```bash
# Build the Go module for Linux (Nakama runs on Linux)
docker run --rm -v $(pwd)/modules:/workspace -w /workspace golang:1.21 bash -c "go mod tidy && go build -buildmode=plugin -o main.so main.go"
```

### 4. Test Matchmaking

#### Single Client Test
```bash
# Test with one client (waits for other players)
npm test
```

#### Multi-Client Test
```bash
# Test with 4 clients automatically
npm run test-multi
```

## 🎮 How It Works

### Matchmaking Flow

1. **Player Connects**: Client authenticates and connects to Nakama
2. **Join Matchmaker**: Player joins matchmaking pool with criteria
3. **Match Found**: Nakama matches players based on criteria
4. **Join Match**: Players automatically join the matched game
5. **Game Play**: Real-time gameplay with score updates

### Match Criteria

The POC uses these matchmaking criteria:
- **Min Players**: 2
- **Max Players**: 4
- **Query**: `"*"` (match any players)
- **String Properties**: `{"mode": "casual"}`
- **Numeric Properties**: `{"skill": 50-150}` (random skill level)

### Match States

- **Waiting**: Players joining, not all ready
- **Playing**: All players ready, game active
- **Finished**: Game completed

## 📁 Project Structure

```
nakama-matchmaking-poc/
├── docker-compose.yml          # Nakama + CockroachDB setup
├── modules/
│   ├── go.mod                  # Go module dependencies
│   ├── main.go                 # Match handler and RPC functions
│   └── main.so                 # Compiled Go module (built)
├── package.json                # Node.js dependencies
├── test-client.js              # Single client test
├── test-multi-client.js        # Multi-client test
└── README.md                   # This file
```

## 🔧 Key Components

### Go Module (`modules/main.go`)

- **MatchmakingMatch**: Implements the match handler
- **Match States**: Tracks players, readiness, scores
- **RPC Functions**: Get match info, player ready, stats
- **Real-time Events**: Player joins/leaves, ready status, score updates

### JavaScript Clients

- **MatchmakingClient**: Handles connection, matchmaking, and gameplay
- **Event Handlers**: Processes match events and updates
- **Game Simulation**: Simulates scoring and gameplay

## 🎯 Matchmaking Features

### Built-in Matchmaker
- Uses Nakama's native matchmaking system
- No custom matchmaking logic required
- Automatic player matching based on criteria

### Real-time Communication
- WebSocket-based real-time messaging
- Match state synchronization
- Player presence tracking

### Match Management
- Automatic match creation
- Player join/leave handling
- Match state transitions

## 🧪 Testing

### Single Client Test
```bash
npm test
```
- Connects one client
- Starts matchmaking
- Waits for other players to join

### Multi-Client Test
```bash
npm run test-multi
```
- Connects 4 clients with delays
- Demonstrates full matchmaking flow
- Shows real-time gameplay simulation

## 📊 Expected Output

When running the multi-client test, you should see:

```
🧪 Testing Nakama Matchmaking POC with Multiple Clients...

🔌 [Player1] Connecting to Nakama...
✅ [Player1] Authenticated - User ID: 12345678-1234-1234-1234-123456789abc
🔗 [Player1] Socket connected
🔍 [Player1] Starting matchmaking...

⏳ [Player2] Waiting 2000ms before connecting...
🔌 [Player2] Connecting to Nakama...
✅ [Player2] Authenticated - User ID: 87654321-4321-4321-4321-cba987654321
🔗 [Player2] Socket connected
🔍 [Player2] Starting matchmaking...

🎯 [Player1] Matchmaker matched! Match ID: match123
🎯 [Player2] Matchmaker matched! Match ID: match123
🎮 [Player1] Joined match: match123
🎮 [Player2] Joined match: match123
✅ [Player1] Marking as ready...
✅ [Player2] Marking as ready...
🚀 [Player1] Match started!
🚀 [Player2] Match started!
📊 [Player1] Current score: 5
📊 [Player2] Current score: 8
```

## 🔍 Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Check what's using the ports
   lsof -i :7350
   lsof -i :26257
   ```

2. **Docker Not Running**
   ```bash
   # Start Docker Desktop, then:
   docker-compose up -d
   ```

3. **Module Not Found**
   ```bash
   # Rebuild the Go module
   docker run --rm -v $(pwd)/modules:/workspace -w /workspace golang:1.21 bash -c "go mod tidy && go build -buildmode=plugin -o main.so main.go"
   ```

### View Logs
```bash
# View Nakama logs
docker-compose logs -f nakama

# View database logs
docker-compose logs -f cockroachdb
```

## 🎯 Next Steps

This POC demonstrates the core matchmaking functionality. You can extend it by:

1. **Custom Match Logic**: Add game-specific rules
2. **Skill-based Matching**: Implement ELO or rating systems
3. **Tournament Support**: Add tournament brackets
4. **Spectator Mode**: Allow non-playing observers
5. **Match Replay**: Record and replay matches

## 📚 Resources

- [Nakama Documentation](https://heroiclabs.com/docs/nakama/)
- [Matchmaker Guide](https://heroiclabs.com/docs/nakama/concepts/matches/)
- [JavaScript Client](https://heroiclabs.com/docs/nakama/client-libraries/javascript-client-guide/)
- [Go Runtime](https://heroiclabs.com/docs/nakama/server-framework/go-runtime/)

---

**Happy Matchmaking! 🎮**
