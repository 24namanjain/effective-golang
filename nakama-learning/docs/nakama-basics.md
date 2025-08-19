# Nakama Basics

A simple guide to understand Nakama's core concepts.

## What is Nakama?

Nakama is a **game server** that handles all the complex backend work for multiplayer games. Instead of building your own server from scratch, Nakama provides:

- **Real-time communication** between players
- **User authentication** and management
- **Leaderboards** and scoring
- **Data storage** for game state
- **Matchmaking** for multiplayer games

## Core Components

### 1. **RPC (Remote Procedure Call)**

Think of RPC as **custom functions** that run on the server. Clients can call these functions to perform actions.

```go
// Server-side: Define a function
nk.RegisterRpc("join_game", func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    // This code runs on the server when a client calls "join_game"
    return `{"success": true, "game_id": "room1"}`, nil
})
```

```javascript
// Client-side: Call the function
const result = await client.rpc(session, "join_game", {player_name: "Alice"});
console.log(result); // {"success": true, "game_id": "room1"}
```

**Common RPC examples:**
- `join_game` - Player joins a game room
- `submit_score` - Player submits their score
- `get_inventory` - Get player's items
- `buy_item` - Purchase something

### 2. **Real-time Communication**

Players can send messages to each other instantly using WebSocket connections.

```javascript
// Join a chat channel
const channel = await socket.joinChat("general", 1, false, false);

// Send a message
await socket.writeChatMessage(channel.id, {message: "Hello everyone!"});

// Listen for messages
socket.onmessage = (message) => {
    console.log("New message:", message);
};
```

### 3. **Leaderboards**

Built-in system for tracking scores and rankings.

```go
// Server: Create a leaderboard
nk.LeaderboardCreate(ctx, "global_scores", "desc", "best", "weekly", 0, false)

// Server: Submit a score
nk.LeaderboardRecordWrite(ctx, "global_scores", userID, username, 1000, nil, nil)
```

```javascript
// Client: Get top scores
const records = await client.listLeaderboardRecords(session, "global_scores", 10);
console.log("Top 10 players:", records.records);
```

### 4. **Storage**

Store player data (inventory, progress, settings, etc.).

```go
// Server: Save player data
objects := []*runtime.StorageWrite{
    &runtime.StorageWrite{
        Collection: "player_data",
        Key:        "inventory",
        UserID:     userID,
        Value:      `{"coins": 100, "items": ["sword", "shield"]}`,
    },
}
nk.StorageWrite(ctx, objects)
```

```javascript
// Client: Read player data
const objects = await client.readStorageObjects(session, [
    {collection: "player_data", key: "inventory"}
]);
console.log("Inventory:", objects.objects[0].value);
```

### 5. **Matches**

Real-time multiplayer game sessions.

```go
// Server: Define match logic
type MyMatch struct {
    players map[string]bool
    scores  map[string]int
}

func (m *MyMatch) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
    m.players = make(map[string]bool)
    m.scores = make(map[string]int)
    return m, nil
}

func (m *MyMatch) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence, metadata map[string]string) (interface{}, []runtime.Presence, error) {
    // Handle player joining
    return state, presences, nil
}
```

## How It All Works Together

### 1. **Client Connects**
```javascript
const client = new nakamajs.Client("defaultkey", "localhost", 7350);
const session = await client.authenticateCustom("user123", "password123");
```

### 2. **Client Calls RPC**
```javascript
// Player wants to join a game
const result = await client.rpc(session, "join_game", {game_type: "battle"});
```

### 3. **Server Processes Request**
```go
// Server receives the RPC call and processes it
// Creates a match, adds player to leaderboard, etc.
```

### 4. **Real-time Updates**
```javascript
// Other players get notified in real-time
socket.onmessage = (message) => {
    if (message.data.type === "player_joined") {
        console.log("New player joined:", message.data.player);
    }
};
```

## Database

Nakama uses **CockroachDB** (default) which is:
- **Distributed**: Can run across multiple servers
- **PostgreSQL compatible**: Uses standard SQL
- **Automatic**: Handles backups, scaling, and failures

## Simple Example: Chat System

Here's how you'd build a simple chat system:

### 1. **Server Setup**
```go
// Create a chat channel
nk.ChannelMessageSend(ctx, "general", "Welcome to the chat!", "system", nil, false)

// Handle chat messages
nk.RegisterRpc("send_message", func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    // Send message to all players in the channel
    nk.ChannelMessageSend(ctx, "general", payload, userID, nil, false)
    return `{"success": true}`, nil
})
```

### 2. **Client Usage**
```javascript
// Join chat channel
const channel = await socket.joinChat("general", 1, false, false);

// Send message
await client.rpc(session, "send_message", "Hello everyone!");

// Listen for messages
socket.onmessage = (message) => {
    console.log("Chat:", message.data.message);
};
```

## Key Benefits

1. **No Server Management**: Nakama handles scaling, backups, and maintenance
2. **Real-time**: Built-in WebSocket support for instant communication
3. **Cross-platform**: Works with any client (Unity, Unreal, Web, Mobile)
4. **Production Ready**: Used by many commercial games
5. **Open Source**: Free to use and modify

## Next Steps

1. **Setup**: Follow the setup guide to get Nakama running
2. **Examples**: Study the code examples to see how everything works
3. **Experiment**: Try modifying the examples to add new features
4. **Build**: Create your own game logic and RPC functions

---

**Remember**: Nakama is just a tool. The game logic, rules, and features are up to you to implement!
