const { Client } = require('@heroiclabs/nakama-js');

// Connect to Nakama server
const client = new Client('defaultkey', 'localhost', 7352);

// Game state
let session = null;
let socket = null;

// Step 1: Connect and login
async function connect() {
    console.log('🔌 Connecting to Nakama...');
    
    // Login (creates account if new)
    session = await client.authenticateCustom('player123', 'password123');
    console.log('✅ Logged in as:', session.username);
    
    // Connect to real-time socket
    socket = client.createSocket();
    await socket.connect(session);
    console.log('✅ Connected to real-time socket');
}

// Step 2: Create a game
async function createGame() {
    console.log('🎮 Creating a new game...');
    
    const result = await client.rpc(session, 'create_game', { name: 'My Game' });
    console.log('✅ Game created:', result);
    return result;
}

// Step 3: Join a game
async function joinGame(gameId) {
    console.log('👥 Joining game:', gameId);
    
    const result = await client.rpc(session, 'join_game', { game_id: gameId });
    console.log('✅ Joined game:', result);
    return result;
}

// Step 4: Submit a score
async function submitScore(score) {
    console.log('🏆 Submitting score:', score);
    
    const result = await client.rpc(session, 'submit_score', { score: score });
    console.log('✅ Score submitted:', result);
    return result;
}

// Step 5: Get leaderboard
async function getLeaderboard() {
    console.log('📊 Getting leaderboard...');
    
    const result = await client.rpc(session, 'get_leaderboard', { limit: 5 });
    console.log('✅ Leaderboard:', result);
    return result;
}

// Step 6: Join chat
async function joinChat() {
    console.log('💬 Joining chat...');
    
    const channel = await socket.joinChat('general', 1, false, false);
    console.log('✅ Joined chat channel');
    
    // Send a message
    await socket.writeChatMessage(channel.id, { message: 'Hello everyone!' });
    console.log('💬 Message sent!');
    
    return channel;
}

// Main game flow
async function playGame() {
    console.log('\n🎮 Starting game...\n');
    
    try {
        // Connect
        await connect();
        
        // Create a game
        const gameResult = await createGame();
        const gameId = gameResult.game_id;
        
        // Submit a random score
        await submitScore(Math.floor(Math.random() * 1000));
        
        // Get leaderboard
        await getLeaderboard();
        
        // Join chat
        await joinChat();
        
        console.log('\n✅ Game completed successfully!\n');
        
    } catch (error) {
        console.error('❌ Game failed:', error.message);
    } finally {
        // Clean up
        if (socket) {
            socket.disconnect();
        }
        console.log('👋 Goodbye!');
    }
}

// Run the example
if (require.main === module) {
    playGame();
}

module.exports = {
    connect,
    createGame,
    joinGame,
    submitScore,
    getLeaderboard,
    joinChat,
    playGame
};
