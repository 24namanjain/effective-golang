// Simple test to verify Nakama setup
const { Client } = require('@heroiclabs/nakama-js');

async function testNakama() {
    console.log('Testing Nakama connection...');
    
    try {
        // Create client
        const client = new Client('defaultkey', '127.0.0.1', 7350, false);
        console.log('Client created successfully');
        
        // Try to authenticate
        const session = await client.authenticateCustom('testuser');
        console.log('Authentication successful:', session.user_id);
        
        // Try to get user account
        const account = await client.getAccount(session);
        console.log('Account retrieved successfully:', account.user.username);
        
        // Try to get user's friends
        const friends = await client.listFriends(session);
        console.log('Friends list retrieved successfully:', friends.friends.length, 'friends found');
        
        console.log('✅ Nakama is working correctly!');
        
    } catch (error) {
        console.error('❌ Nakama test failed:', error.message);
    }
}

testNakama();
