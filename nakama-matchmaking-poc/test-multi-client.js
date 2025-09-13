// Nakama Matchmaking POC - Multi-Client Test
const { Client } = require('@heroiclabs/nakama-js');

class MatchmakingClient {
    constructor(username, delay = 0) {
        this.username = username;
        this.client = new Client('defaultkey', '127.0.0.1', 7350, false);
        this.socket = null;
        this.session = null;
        this.currentMatch = null;
        this.delay = delay;
    }

    async connect() {
        try {
            if (this.delay > 0) {
                console.log(`⏳ [${this.username}] Waiting ${this.delay}ms before connecting...`);
                await new Promise(resolve => setTimeout(resolve, this.delay));
            }
            
            console.log(`🔌 [${this.username}] Connecting to Nakama...`);
            
            // Authenticate
            this.session = await this.client.authenticateCustom(this.username);
            console.log(`✅ [${this.username}] Authenticated - User ID: ${this.session.user_id}`);
            
            // Create socket connection
            this.socket = this.client.createSocket();
            await this.socket.connect(this.session);
            console.log(`🔗 [${this.username}] Socket connected`);
            
            // Set up event handlers
            this.setupEventHandlers();
            
        } catch (error) {
            console.error(`❌ [${this.username}] Connection failed:`, error.message);
            throw error;
        }
    }

    setupEventHandlers() {
        // Handle matchmaker matched event
        this.socket.onmatchmakermatched = async (matched) => {
            console.log(`🎯 [${this.username}] Matchmaker matched! Match ID: ${matched.match_id}`);
            console.log(`📊 [${this.username}] Match details:`, {
                matchId: matched.match_id,
                token: matched.token,
                users: matched.users.map(u => ({ id: u.user_id, username: u.username }))
            });
            
            // Join the match
            this.currentMatch = await this.socket.joinMatch(matched.match_id, matched.token);
            console.log(`🎮 [${this.username}] Joined match: ${this.currentMatch.match_id}`);
            
            // Auto-mark as ready after a short delay
            setTimeout(() => {
                this.markReady();
            }, 1000);
        };

        // Handle match data (messages from other players)
        this.socket.onmatchdata = (matchData) => {
            try {
                const data = JSON.parse(matchData.data);
                console.log(`📨 [${this.username}] Received match data:`, data);
            
            switch (data.type) {
                case 'player_joined':
                    console.log(`👋 [${this.username}] Player joined: ${data.player.username}`);
                    break;
                case 'player_left':
                    console.log(`👋 [${this.username}] Player left: ${data.username}`);
                    break;
                case 'player_ready':
                    console.log(`✅ [${this.username}] Player ready: ${data.username}`);
                    break;
                case 'match_started':
                    console.log(`🚀 [${this.username}] Match started!`);
                    this.simulateGameplay();
                    break;
                case 'score_update':
                    console.log(`📊 [${this.username}] Score update: ${data.username} = ${data.score}`);
                    break;
                case 'match_ending':
                    console.log(`🏁 [${this.username}] Match ending...`);
                    break;
            }
            } catch (error) {
                console.log(`⚠️ [${this.username}] Failed to parse match data:`, matchData.data, error.message);
            }
        };

        // Handle match presence (players joining/leaving)
        this.socket.onmatchpresence = (matchPresence) => {
            console.log(`👥 [${this.username}] Match presence update:`, {
                joins: matchPresence.joins ? matchPresence.joins.map(p => ({ id: p.user_id, username: p.username })) : [],
                leaves: matchPresence.leaves ? matchPresence.leaves.map(p => ({ id: p.user_id, username: p.username })) : []
            });
        };
    }

    async startMatchmaking() {
        try {
            console.log(`🔍 [${this.username}] Starting matchmaking...`);
            
            // Add to matchmaker with skill-based matching
            const minPlayers = 2;
            const maxPlayers = 4;
            const query = "*"; // Match any players
            const stringProperties = { "mode": "casual" };
            const numericProperties = { "skill": Math.floor(Math.random() * 100) + 50 }; // Random skill 50-150
            
            const matchmakerTicket = await this.socket.addMatchmaker(
                query,
                minPlayers,
                maxPlayers,
                stringProperties,
                numericProperties
            );
            
            console.log(`🎫 [${this.username}] Matchmaker ticket: ${matchmakerTicket.ticket}`);
            console.log(`📋 [${this.username}] Matchmaking criteria:`, {
                minPlayers,
                maxPlayers,
                query,
                stringProperties,
                numericProperties
            });
            
        } catch (error) {
            console.error(`❌ [${this.username}] Matchmaking failed:`, error.message);
            throw error;
        }
    }

    async markReady() {
        if (!this.currentMatch) {
            console.log(`⚠️ [${this.username}] No active match to mark ready`);
            return;
        }

        try {
            console.log(`✅ [${this.username}] Marking as ready...`);
            
            // Send ready message to match
            const readyMessage = JSON.stringify({ type: "player_ready" });
            await this.socket.sendMatchState(this.currentMatch.match_id, 1, readyMessage);
            
        } catch (error) {
            console.error(`❌ [${this.username}] Failed to mark ready:`, error.message);
        }
    }

    async simulateGameplay() {
        if (!this.currentMatch) return;

        console.log(`🎮 [${this.username}] Starting gameplay simulation...`);
        
        // Simulate scoring
        let score = 0;
        const scoreInterval = setInterval(async () => {
            if (!this.currentMatch) {
                clearInterval(scoreInterval);
                return;
            }

            score += Math.floor(Math.random() * 10) + 1;
            console.log(`📊 [${this.username}] Current score: ${score}`);
            
            try {
                const scoreMessage = JSON.stringify({ 
                    type: "player_score", 
                    score: score 
                });
                await this.socket.sendMatchState(this.currentMatch.match_id, 1, scoreMessage);
            } catch (error) {
                console.error(`❌ [${this.username}] Failed to send score:`, error.message);
                clearInterval(scoreInterval);
            }
        }, 2000);

        // Stop after 20 seconds
        setTimeout(() => {
            clearInterval(scoreInterval);
            console.log(`🏁 [${this.username}] Gameplay simulation ended`);
        }, 20000);
    }

    async disconnect() {
        try {
            if (this.socket) {
                await this.socket.disconnect();
                console.log(`🔌 [${this.username}] Disconnected`);
            }
        } catch (error) {
            console.error(`❌ [${this.username}] Disconnect error:`, error.message);
        }
    }
}

// Multi-client test function
async function testMultiClientMatchmaking() {
    console.log('🧪 Testing Nakama Matchmaking POC with Multiple Clients...\n');
    
    const clients = [
        new MatchmakingClient('Player1', 0),
        new MatchmakingClient('Player2', 2000),
        new MatchmakingClient('Player3', 4000),
        new MatchmakingClient('Player4', 6000)
    ];
    
    try {
        // Connect all clients
        for (const client of clients) {
            await client.connect();
        }
        
        console.log('\n🔍 Starting matchmaking for all clients...\n');
        
        // Start matchmaking for all clients
        for (const client of clients) {
            await client.startMatchmaking();
        }
        
        console.log('\n⏳ All clients are now matchmaking... (Press Ctrl+C to stop)');
        console.log('💡 Watch as players get matched and games start!\n');
        
        // Keep the process alive
        process.on('SIGINT', async () => {
            console.log('\n🛑 Shutting down all clients...');
            for (const client of clients) {
                await client.disconnect();
            }
            process.exit(0);
        });
        
    } catch (error) {
        console.error('❌ Multi-client test failed:', error.message);
        console.log('\n💡 Make sure Nakama is running: docker-compose up -d');
        
        // Clean up on error
        for (const client of clients) {
            await client.disconnect();
        }
        process.exit(1);
    }
}

// Run the multi-client test
testMultiClientMatchmaking();
