- Nakama handles the backend logic which is not related to the actual gameplay
  - It is also responsible for the communication between the players
  - It is also responsible for the storage of the game data
  - It is also responsible for the leaderboards
  - It is also responsible for the matchmaking
  - It is also responsible for the realtime socket
  - It is also responsible for the social features

Nakama Building Blocks

1. RPCs
2. Matches
    - Matches are the actual game sessions
    - They are the main way to handle the game logic
    - Handle games state, player interactions, and sync
    - can be used for turn-based or real-time
3. Matchmaker
    - Matchmaker is the system that handles the matchmaking of the players
    - supports various criteria like skill level, game mode, etc.
    - can create teams, handle queue management and more
4. Storage
    - Storage is the system that handles the storage of the game data
    - can store game state, player data, settings, etc.
    - can be used for persistent or non-persistent data
5. Leaderboards
    - Leaderboards are the system that handles the leaderboards of the game
    - can store and retrieve leaderboard data
    - can be used for global or regional leaderboards
6. Social
    - Social is the system that handles the social features of the game
    - can handle friend requests, chat, groups, etc.
7. Realtime Socket
    - Realtime Socket is the system that handles the realtime socket of the game
    - can handle realtime communication between the players
    - can be used for chat, match data streams, and presence events