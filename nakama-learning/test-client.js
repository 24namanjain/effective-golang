// Simple test to check if Nakama is working
const { Client } = require('@heroiclabs/nakama-js');

async function testNakama() {
    console.log('🧪 Testing Nakama connection...');
    
    try {
        // Create client and connect
        const client = new Client('defaultkey', '127.0.0.1', 7352, false);
        console.log('✅ Client created');
        
        // Try to login
        const session = await client.authenticateCustom('testuser');
        console.log('✅ Login successful - User ID:', session.user_id);
        
        // Get user account info
        const account = await client.getAccount(session);
        console.log('✅ Account retrieved - Username:', account.user.username);
        
        // Get friends list (should be empty for new user)
        const friends = await client.listFriends(session);
        console.log('✅ Friends list retrieved - Count:', friends.friends.length);
        
        console.log('\n🎉 All tests passed! Nakama is working correctly.\n');
        
    } catch (error) {
        console.error('❌ Test failed:', error.message);
        console.log('\n💡 Make sure Nakama is running: docker compose up -d');
    }
}

testNakama();
