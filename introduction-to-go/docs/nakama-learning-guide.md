# Nakama Learning Guide for Go Developers

## What is Nakama?

Nakama is an open-source game server that provides:
- **Real-time multiplayer** functionality
- **User authentication** and social features
- **Leaderboards** and achievements
- **Matchmaking** and game sessions
- **Storage** for user data
- **RPC** for custom server logic

Think of it as a "backend-as-a-service" specifically designed for games, similar to Firebase but with game-specific features.

## Architecture Overview

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │    │   Client    │    │   Client    │
│  (Mobile/   │    │  (Mobile/   │    │  (Mobile/   │
│   Web/PC)   │    │   Web/PC)   │    │   Web/PC)   │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                    ┌─────────────┐
                    │   Nakama    │
                    │   Server    │
                    │  (Go/Lua)   │
                    └─────────────┘
                           │
                    ┌─────────────┐
                    │  Database   │
                    │ (PostgreSQL)│
                    └─────────────┘
```

## Core Concepts

### 1. **RPC (Remote Procedure Call)**
Custom server functions that clients can call. Think of them as API endpoints.

```go
// Example: Custom RPC function
func JoinGameRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    // Parse the request
    var request struct {
        GameID string `json:"game_id"`
    }
    if err := json.Unmarshal([]byte(payload), &request); err != nil {
        return "", err
    }
    
    // Get the user ID from context
    userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
    if !ok {
        return "", errors.New("user not authenticated")
    }
    
    // Join the game logic here
    // ...
    
    // Return success response
    response := map[string]interface{}{
        "success": true,
        "game_id": request.GameID,
    }
    
    responseBytes, _ := json.Marshal(response)
    return string(responseBytes), nil
}

// Register the RPC
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
    initializer.RegisterRpc("join_game", JoinGameRPC)
    return nil
}
```

### 2. **Events**
System events that trigger your custom code when certain actions happen.

```go
// Example: Track login streaks
func AfterAuthenticateHook(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, out *api.Session, in *api.AuthenticateRequest) error {
    userID := out.UserId
    
    // Get user's last login
    objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{
        {
            Collection: "user_stats",
            Key:        "login_streak",
            UserID:     userID,
        },
    })
    
    if err != nil {
        logger.Error("Failed to read user stats: %v", err)
        return err
    }
    
    var streak int
    if len(objects) > 0 {
        var stats struct {
            Streak int `json:"streak"`
            LastLogin time.Time `json:"last_login"`
        }
        json.Unmarshal([]byte(objects[0].Value), &stats)
        streak = stats.Streak
    }
    
    // Update streak
    now := time.Now()
    newStats := map[string]interface{}{
        "streak": streak + 1,
        "last_login": now,
    }
    
    statsBytes, _ := json.Marshal(newStats)
    _, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{
        {
            Collection: "user_stats",
            Key:        "login_streak",
            UserID:     userID,
            Value:      string(statsBytes),
        },
    })
    
    return err
}
```

### 3. **Leaderboards**
Competitive scoring system with rankings.

```go
// Example: Create and manage leaderboards
func CreateLeaderboards(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) error {
    // Create a global leaderboard
    _, err := nk.LeaderboardCreate(ctx, "global_scores", "Global High Scores", "desc", "best", "0 0 * * *", false, "")
    if err != nil {
        return err
    }
    
    // Create a weekly leaderboard (resets every Monday)
    _, err = nk.LeaderboardCreate(ctx, "weekly_scores", "Weekly Scores", "desc", "best", "0 0 * * 1", false, "")
    if err != nil {
        return err
    }
    
    return nil
}

// Submit a score
func SubmitScoreRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    var request struct {
        Score int64 `json:"score"`
    }
    if err := json.Unmarshal([]byte(payload), &request); err != nil {
        return "", err
    }
    
    userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
    
    // Submit to global leaderboard
    _, err := nk.LeaderboardRecordWrite(ctx, "global_scores", userID, request.Score, nil, nil)
    if err != nil {
        return "", err
    }
    
    // Submit to weekly leaderboard
    _, err = nk.LeaderboardRecordWrite(ctx, "weekly_scores", userID, request.Score, nil, nil)
    if err != nil {
        return "", err
    }
    
    return `{"success": true}`, nil
}

// Get leaderboard rankings
func GetLeaderboardRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    var request struct {
        LeaderboardID string `json:"leaderboard_id"`
        Limit         int    `json:"limit"`
    }
    if err := json.Unmarshal([]byte(payload), &request); err != nil {
        return "", err
    }
    
    if request.Limit == 0 {
        request.Limit = 10
    }
    
    records, err := nk.LeaderboardRecordsList(ctx, request.LeaderboardID, []string{}, request.Limit, "")
    if err != nil {
        return "", err
    }
    
    response := map[string]interface{}{
        "records": records.Records,
        "owner_records": records.OwnerRecords,
    }
    
    responseBytes, _ := json.Marshal(response)
    return string(responseBytes), nil
}
```

### 4. **Matches & Matchmaker**
Real-time multiplayer game sessions.

```go
// Example: Simple match implementation
type MatchState struct {
    Players map[string]*Player
    GameState string
    StartTime time.Time
}

type Player struct {
    UserID string
    Score  int
    Ready  bool
}

func MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
    return &Match{
        state: &MatchState{
            Players: make(map[string]*Player),
            GameState: "waiting",
        },
    }, nil
}

type Match struct {
    state *MatchState
}

func (m *Match) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
    return m, nil
}

func (m *Match) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
    if len(m.state.Players) >= 4 {
        return m.state, false, "Match is full"
    }
    
    m.state.Players[presence.UserId] = &Player{
        UserID: presence.UserId,
        Score:  0,
        Ready:  false,
    }
    
    return m.state, true, ""
}

func (m *Match) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
    for _, presence := range presences {
        logger.Info("Player joined: %s", presence.UserId)
    }
    
    // Start game if enough players
    if len(m.state.Players) >= 2 {
        m.state.GameState = "playing"
        m.state.StartTime = time.Now()
        
        // Notify all players
        dispatcher.BroadcastMessage(1, &api.Envelope{
            Message: &api.Envelope_GameStart{
                GameStart: &api.GameStart{
                    StartTime: m.state.StartTime.Unix(),
                },
            },
        }, nil, nil, true)
    }
    
    return m.state
}

func (m *Match) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
    for _, presence := range presences {
        delete(m.state.Players, presence.UserId)
        logger.Info("Player left: %s", presence.UserId)
    }
    
    return m.state
}

func (m *Match) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
    // Handle incoming messages
    for _, message := range messages {
        switch message.OpCode {
        case 1: // Player ready
            if player, exists := m.state.Players[message.UserId]; exists {
                player.Ready = true
            }
        case 2: // Update score
            var scoreUpdate struct {
                Score int `json:"score"`
            }
            if err := json.Unmarshal(message.Data, &scoreUpdate); err == nil {
                if player, exists := m.state.Players[message.UserId]; exists {
                    player.Score = scoreUpdate.Score
                }
            }
        }
    }
    
    return m.state
}

func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
    // Save final scores to leaderboard
    for userID, player := range m.state.Players {
        nk.LeaderboardRecordWrite(ctx, "global_scores", userID, int64(player.Score), nil, nil)
    }
    
    return m.state
}
```

### 5. **Storage Engine**
Key-value storage for user data.

```go
// Example: User profile management
func SaveUserProfileRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    var profile struct {
        DisplayName string `json:"display_name"`
        Avatar      string `json:"avatar"`
        Level       int    `json:"level"`
        Experience  int    `json:"experience"`
    }
    
    if err := json.Unmarshal([]byte(payload), &profile); err != nil {
        return "", err
    }
    
    userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
    
    // Save profile
    profileBytes, _ := json.Marshal(profile)
    _, err := nk.StorageWrite(ctx, []*runtime.StorageWrite{
        {
            Collection: "user_profiles",
            Key:        "profile",
            UserID:     userID,
            Value:      string(profileBytes),
        },
    })
    
    if err != nil {
        return "", err
    }
    
    return `{"success": true}`, nil
}

func GetUserProfileRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    var request struct {
        UserID string `json:"user_id"`
    }
    
    if err := json.Unmarshal([]byte(payload), &request); err != nil {
        return "", err
    }
    
    // If no user_id provided, get current user's profile
    if request.UserID == "" {
        request.UserID, _ = ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
    }
    
    objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{
        {
            Collection: "user_profiles",
            Key:        "profile",
            UserID:     request.UserID,
        },
    })
    
    if err != nil {
        return "", err
    }
    
    if len(objects) == 0 {
        return `{"error": "Profile not found"}`, nil
    }
    
    return objects[0].Value, nil
}
```

### 6. **Realtime Events**
Handle real-time communication.

```go
// Example: Chat system
func ChatMessageRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    var message struct {
        Channel string `json:"channel"`
        Message string `json:"message"`
    }
    
    if err := json.Unmarshal([]byte(payload), &message); err != nil {
        return "", err
    }
    
    userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
    
    // Join the channel if not already joined
    _, err := nk.ChannelJoin(ctx, userID, message.Channel, false, false, "")
    if err != nil {
        return "", err
    }
    
    // Send the message
    _, err = nk.ChannelMessageSend(ctx, message.Channel, &api.Envelope{
        Message: &api.Envelope_ChatMessage{
            ChatMessage: &api.ChatMessage{
                UserId:  userID,
                Message: message.Message,
            },
        },
    })
    
    if err != nil {
        return "", err
    }
    
    return `{"success": true}`, nil
}
```

## Getting Started Steps

### 1. **Install Nakama**
```bash
# Using Docker (recommended)
docker run -d --name nakama -p 7350:7350 -p 7351:7351 heroiclabs/nakama:latest

# Or download binary from GitHub
```

### 2. **Create Your First Module**
Create a file `main.go`:

```go
package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "errors"
    
    "github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
    // Register RPC functions
    initializer.RegisterRpc("hello_world", HelloWorldRPC)
    initializer.RegisterRpc("submit_score", SubmitScoreRPC)
    
    // Register hooks
    initializer.RegisterAfterAuthenticate(AfterAuthenticateHook)
    
    // Create leaderboards
    if err := CreateLeaderboards(ctx, logger, db, nk); err != nil {
        return err
    }
    
    logger.Info("Module initialized successfully!")
    return nil
}

func HelloWorldRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
    
    response := map[string]interface{}{
        "message": "Hello from Nakama!",
        "user_id": userID,
    }
    
    responseBytes, _ := json.Marshal(response)
    return string(responseBytes), nil
}
```

### 3. **Test Your Module**
```bash
# Build and run with your module
nakama --runtime.path /path/to/your/module
```

### 4. **Client Integration**
```javascript
// JavaScript client example
const client = new nakamajs.Client("defaultkey", "127.0.0.1", 7350, false);

// Authenticate
const session = await client.authenticateCustom("user123");

// Call your RPC
const result = await client.rpc(session, "hello_world", {});
console.log(result.payload);
```

## Best Practices

1. **Error Handling**: Always return meaningful error messages
2. **Validation**: Validate all input data
3. **Logging**: Use logger for debugging and monitoring
4. **Performance**: Use batch operations for multiple database calls
5. **Security**: Validate user permissions in RPC functions
6. **Testing**: Test your modules thoroughly

## Common Patterns

### Authentication Flow
```go
func AuthenticateUser(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    var request struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    
    if err := json.Unmarshal([]byte(payload), &request); err != nil {
        return "", err
    }
    
    // Validate credentials
    // Create or get user account
    // Return session token
    
    return `{"token": "session_token"}`, nil
}
```

### Data Validation
```go
func ValidateUserInput(input map[string]interface{}) error {
    if username, ok := input["username"].(string); ok {
        if len(username) < 3 || len(username) > 20 {
            return errors.New("username must be 3-20 characters")
        }
    }
    return nil
}
```

This guide should give you a solid foundation to start building with Nakama! Start with simple RPC functions and gradually add more complex features like matches and real-time events.
