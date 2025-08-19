// =============================================================================
// NAKAMA CLIENT EXAMPLE - JAVASCRIPT
// =============================================================================
// This example demonstrates how to build a client that interacts with the Nakama game server
// It shows authentication, game creation, joining, and real-time gameplay
// 
// Key Concepts:
// - Client-Server Communication: How clients talk to the Nakama server
// - RPC Calls: Remote procedure calls to server functions
// - Session Management: Handling user authentication
// - Game State Management: Tracking and updating game data
// - Real-time Features: Chat, presence, and live updates
// =============================================================================

const { Client } = require('@heroiclabs/nakama-js');

// =============================================================================
// CLIENT INITIALIZATION
// =============================================================================

// Initialize the Nakama client
// Parameters: server key, host, port, use SSL
const client = new Client('defaultkey', '127.0.0.1', 7350, false);

// Global state variables to track current session and game
let session = null;        // User's authentication session
let currentGame = null;    // Current game the user is in

// =============================================================================
// AUTHENTICATION FUNCTIONS
// =============================================================================

// authenticateUser - Log in a user to the game server
// This is the first step before any game operations
// Parameters: username (string) - The user's display name
// Returns: session object with user information
async function authenticateUser(username) {
    try {
        // Authenticate with custom ID (in production, use proper auth like Google, Facebook, etc.)
        // This creates a session that identifies the user to the server
        session = await client.authenticateCustom(username);
        console.log(`✅ Authenticated as: ${username}`);
        console.log(`📋 User ID: ${session.user_id}`);
        return session;
    } catch (error) {
        console.error('❌ Authentication failed:', error);
        throw error;
    }
}

// =============================================================================
// GAME MANAGEMENT FUNCTIONS
// =============================================================================

// createGame - Create a new multiplayer game session
// This calls the server's CreateGameRPC function
// Parameters: maxPlayers (number) - Maximum players allowed, minPlayers (number) - Minimum players needed
// Returns: game object with game state
async function createGame(maxPlayers = 4, minPlayers = 2) {
    try {
        // Prepare the request payload (this becomes the 'payload' parameter in the server RPC)
        const payload = {
            max_players: maxPlayers,  // Server will validate this (2-8 players)
            min_players: minPlayers   // Minimum players needed to start
        };
        
        // Call the server's 'create_game' RPC function
        // This sends the payload to the server and waits for a response
        const result = await client.rpc(session, 'create_game', payload);
        
        // Parse the JSON response from the server
        const response = JSON.parse(result.payload);
        
        if (response.success) {
            // Store the game state locally
            currentGame = response.game;
            console.log(`🎮 Game created: ${response.game_id}`);
            console.log('📊 Game state:', currentGame);
            return currentGame;
        } else {
            throw new Error('Failed to create game');
        }
    } catch (error) {
        console.error('❌ Create game failed:', error);
        throw error;
    }
}

// joinGame - Join an existing game created by another player
// This calls the server's JoinGameRPC function
// Parameters: gameId (string) - The ID of the game to join
// Returns: updated game object
async function joinGame(gameId) {
    try {
        // Prepare the request payload
        const payload = {
            game_id: gameId  // The unique identifier of the game
        };
        
        // Call the server's 'join_game' RPC function
        const result = await client.rpc(session, 'join_game', payload);
        const response = JSON.parse(result.payload);
        
        if (response.success) {
            // Update local game state with server response
            currentGame = response.game;
            console.log(`🎯 Joined game: ${gameId}`);
            console.log('📊 Game state:', currentGame);
            return currentGame;
        } else {
            throw new Error('Failed to join game');
        }
    } catch (error) {
        console.error('❌ Join game failed:', error);
        throw error;
    }
}

// startGame - Start the actual gameplay when all players are ready
// This calls the server's StartGameRPC function
// Returns: updated game object
async function startGame() {
    if (!currentGame) {
        throw new Error('No current game');
    }
    
    try {
        const payload = {
            game_id: currentGame.game_id
        };
        
        // Call the server's 'start_game' RPC function
        const result = await client.rpc(session, 'start_game', payload);
        const response = JSON.parse(result.payload);
        
        if (response.success) {
            currentGame = response.game;
            console.log('🚀 Game started!');
            console.log('📊 Game state:', currentGame);
            return currentGame;
        } else {
            throw new Error('Failed to start game');
        }
    } catch (error) {
        console.error('❌ Start game failed:', error);
        throw error;
    }
}

// =============================================================================
// GAMEPLAY FUNCTIONS
// =============================================================================

// updatePosition - Update player position during gameplay
// This is called frequently as the player moves around
// Parameters: x (number), y (number) - New coordinates
// Returns: updated position
async function updatePosition(x, y) {
    if (!currentGame) {
        throw new Error('No current game');
    }
    
    try {
        const payload = {
            game_id: currentGame.game_id,
            position: {
                x: x,  // X coordinate (horizontal position)
                y: y   // Y coordinate (vertical position)
            }
        };
        
        // Call the server's 'update_position' RPC function
        const result = await client.rpc(session, 'update_position', payload);
        const response = JSON.parse(result.payload);
        
        if (response.success) {
            console.log(`📍 Position updated: (${x}, ${y})`);
            return response.position;
        } else {
            throw new Error('Failed to update position');
        }
    } catch (error) {
        console.error('❌ Update position failed:', error);
        throw error;
    }
}

// endGame - End the current game and save scores
// This calls the server's EndGameRPC function
// Returns: final game state
async function endGame() {
    if (!currentGame) {
        throw new Error('No current game');
    }
    
    try {
        const payload = {
            game_id: currentGame.game_id
        };
        
        // Call the server's 'end_game' RPC function
        const result = await client.rpc(session, 'end_game', payload);
        const response = JSON.parse(result.payload);
        
        if (response.success) {
            currentGame = response.game;
            console.log('🏁 Game ended!');
            console.log('📊 Final game state:', currentGame);
            return currentGame;
        } else {
            throw new Error('Failed to end game');
        }
    } catch (error) {
        console.error('❌ End game failed:', error);
        throw error;
    }
}

// =============================================================================
// GAME STATE FUNCTIONS
// =============================================================================

// getGameState - Get the current state of a game
// Useful for syncing client state with server
// Parameters: gameId (string) - The game ID
// Returns: current game state
async function getGameState(gameId) {
    try {
        const payload = {
            game_id: gameId
        };
        
        // Call the server's 'get_game_state' RPC function
        const result = await client.rpc(session, 'get_game_state', payload);
        const response = JSON.parse(result.payload);
        
        if (response.success) {
            console.log('📊 Retrieved game state');
            return JSON.parse(response.game); // Parse the game JSON string
        } else {
            throw new Error('Failed to get game state');
        }
    } catch (error) {
        console.error('❌ Get game state failed:', error);
        throw error;
    }
}

// getLeaderboard - Get the global leaderboard
// Shows competitive rankings across all games
// Parameters: limit (number) - Number of records to return
// Returns: leaderboard records
async function getLeaderboard(limit = 10) {
    try {
        const payload = {
            limit: limit
        };
        
        // Call the server's 'get_leaderboard' RPC function
        const result = await client.rpc(session, 'get_leaderboard', payload);
        const response = JSON.parse(result.payload);
        
        if (response.success) {
            console.log('🏆 Retrieved leaderboard');
            console.log('📊 Records:', response.records);
            return response.records;
        } else {
            throw new Error('Failed to get leaderboard');
        }
    } catch (error) {
        console.error('❌ Get leaderboard failed:', error);
        throw error;
    }
}

// =============================================================================
// REAL-TIME FEATURES
// =============================================================================

// joinChatChannel - Join a chat channel for real-time communication
// Parameters: channelName (string) - Name of the channel to join
// Returns: channel information
async function joinChatChannel(channelName) {
    try {
        // Join a chat channel (room for real-time messaging)
        const channel = await client.joinChat(session, channelName, 1, false, false);
        console.log(`💬 Joined chat channel: ${channelName}`);
        return channel;
    } catch (error) {
        console.error('❌ Join chat failed:', error);
        throw error;
    }
}

// sendChatMessage - Send a message to a chat channel
// Parameters: channelId (string), message (string) - The message to send
async function sendChatMessage(channelId, message) {
    try {
        // Send a message to the specified channel
        const messageAck = await client.writeChatMessage(session, channelId, {
            content: message
        });
        console.log(`💬 Message sent: ${message}`);
        return messageAck;
    } catch (error) {
        console.error('❌ Send message failed:', error);
        throw error;
    }
}

// =============================================================================
// REAL-TIME CONNECTION SETUP
// =============================================================================

// setupRealtimeConnection - Establish real-time connection for live updates
// This enables features like live chat, presence updates, and real-time game state
async function setupRealtimeConnection() {
    try {
        // Create a real-time socket connection
        const socket = client.createSocket();
        
        // Connect to the server
        await socket.connect(session);
        console.log('🔌 Real-time connection established');
        
        // Set up event listeners for real-time events
        
        // Listen for chat messages
        socket.onmessage = (message) => {
            console.log('💬 Received message:', message);
        };
        
        // Listen for presence updates (players joining/leaving)
        socket.onpresence = (presence) => {
            console.log('👥 Presence update:', presence);
        };
        
        // Listen for channel join events
        socket.onchannelmessage = (message) => {
            console.log('📢 Channel message:', message);
        };
        
        return socket;
    } catch (error) {
        console.error('❌ Real-time connection failed:', error);
        throw error;
    }
}

// =============================================================================
// EXAMPLE USAGE - COMPLETE GAME FLOW
// =============================================================================

// runExample - Demonstrates a complete game flow
// This shows how all the functions work together
async function runExample() {
    try {
        console.log('🎮 Starting Nakama Game Example...\n');
        
        // Step 1: Authenticate user
        console.log('1️⃣ Authenticating user...');
        await authenticateUser('player1');
        
        // Step 2: Create a new game
        console.log('\n2️⃣ Creating a new game...');
        await createGame(4, 2);
        
        // Step 3: Get leaderboard
        console.log('\n3️⃣ Getting leaderboard...');
        await getLeaderboard(5);
        
        // Step 4: Join chat channel
        console.log('\n4️⃣ Joining chat channel...');
        const channel = await joinChatChannel('general');
        
        // Step 5: Send a chat message
        console.log('\n5️⃣ Sending chat message...');
        await sendChatMessage(channel.id, 'Hello, everyone!');
        
        // Step 6: Set up real-time connection
        console.log('\n6️⃣ Setting up real-time connection...');
        const socket = await setupRealtimeConnection();
        
        // Step 7: Simulate gameplay
        console.log('\n7️⃣ Simulating gameplay...');
        await updatePosition(10, 20);
        await updatePosition(15, 25);
        
        // Step 8: End the game
        console.log('\n8️⃣ Ending the game...');
        await endGame();
        
        console.log('\n✅ Example completed successfully!');
        
    } catch (error) {
        console.error('❌ Example failed:', error);
    }
}

// =============================================================================
// EXPORT FUNCTIONS FOR USE IN OTHER MODULES
// =============================================================================

module.exports = {
    authenticateUser,
    createGame,
    joinGame,
    startGame,
    updatePosition,
    endGame,
    getGameState,
    getLeaderboard,
    joinChatChannel,
    sendChatMessage,
    setupRealtimeConnection,
    runExample
};

// =============================================================================
// USAGE INSTRUCTIONS
// =============================================================================
/*
To use this client:

1. Make sure Nakama server is running (see setup guide)
2. Install dependencies: npm install @heroiclabs/nakama-js
3. Run the example: node nakama_client_example.js

Or import and use individual functions:

const { authenticateUser, createGame } = require('./nakama_client_example.js');

async function myGame() {
    await authenticateUser('myplayer');
    await createGame(4, 2);
}
*/
