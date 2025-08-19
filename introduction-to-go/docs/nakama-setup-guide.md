# Nakama Quick Setup Guide

## Prerequisites

- Go 1.19 or later
- Docker (recommended) 
- Node.js (for client examples)

## Step 1: Install Nakama with CockroachDB

### Using Docker Compose (Recommended)
```bash
# Create docker-compose.yml
services:
  cockroachdb:
    image: cockroachdb/cockroach:latest-v24.1
    command: start-single-node --insecure --store=attrs=ssd,path=/var/lib/cockroach/
    restart: "no"
    volumes:
      - data:/var/lib/cockroach
    expose:
      - "8080"
      - "26257"
    ports:
      - "26257:26257"
      - "8080:8080"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health?ready=1"]
      interval: 3s
      timeout: 3s
      retries: 5

  nakama:
    image: registry.heroiclabs.com/heroiclabs/nakama:3.26.0
    entrypoint:
      - "/bin/sh"
      - "-ecx"
      - |
        /nakama/nakama migrate up --database.address root@cockroachdb:26257 &&
        exec /nakama/nakama --name nakama1 --database.address root@cockroachdb:26257 --logger.level DEBUG --session.token_expiry_sec 7200 --metrics.prometheus_port 9100
    restart: "no"
    links:
      - "cockroachdb:db"
    depends_on:
      cockroachdb:
        condition: service_healthy
    volumes:
      - ./data:/nakama/data
      - ./modules:/nakama/modules
    expose:
      - "7349"
      - "7350"
      - "7351"
      - "9100"
    ports:
      - "7349:7349"
      - "7350:7350"
      - "7351:7351"

volumes:
  data:

# Start the services
docker compose up -d
```

### Alternative: Using Binary
```bash
# Download from GitHub releases
wget https://github.com/heroiclabs/nakama/releases/latest/download/nakama.tar.gz
tar -xzf nakama.tar.gz
./nakama
```

## Step 2: Create Your Go Module

```bash
mkdir modules
cd modules
go mod init nakama-module
```

Create `modules/main.go`:
```go
package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "time"

    "github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
    // Register RPC functions
    initializer.RegisterRpc("hello_world", HelloWorldRPC)
    initializer.RegisterRpc("get_time", GetTimeRPC)
    
    logger.Info("Simple Nakama Module initialized successfully!")
    return nil
}

func HelloWorldRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
    
    response := map[string]interface{}{
        "message": "Hello from Nakama!",
        "user_id": userID,
        "timestamp": time.Now().Unix(),
    }
    
    responseBytes, _ := json.Marshal(response)
    return string(responseBytes), nil
}

func GetTimeRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
    response := map[string]interface{}{
        "server_time": time.Now().Format(time.RFC3339),
        "unix_time": time.Now().Unix(),
    }
    
    responseBytes, _ := json.Marshal(response)
    return string(responseBytes), nil
}
```

Create `modules/go.mod`:
```go
module nakama-module

go 1.21

require github.com/heroiclabs/nakama-common v1.32.0
```

## Step 3: Build and Deploy Module

### Build for Linux (required for Docker container)
```bash
# Build the module using Docker (cross-platform)
docker run --rm -v $(pwd)/modules:/app -w /app golang:1.21 sh -c "go mod tidy && go build -buildmode=plugin -o main.so main.go"

# Copy to data directory
cp modules/main.so data/modules/

# Restart Nakama to load the module
docker compose restart nakama
```

### Alternative: Build locally (for development)
```bash
cd modules
go mod tidy
go build -buildmode=plugin -o main.so main.go
```

## Step 4: Test Your Setup

### Using the Nakama Console
1. Open http://localhost:7351 in your browser
2. Use the default credentials:
   - Username: `admin`
   - Password: `password`

### Using the JavaScript Client
```bash
npm install @heroiclabs/nakama-js
```

Create `test-nakama.js`:
```javascript
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
```

Run the test:
```bash
node test-nakama.js
```

### Using HTTP API (curl)
```bash
# Authenticate a user
curl -s -X POST "http://127.0.0.1:7350/v2/account/authenticate/custom?create=true" \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic ZGVmYXVsdGtleTo=" \
  -d '{"id":"user123"}' | jq .

# Get account info (replace TOKEN with the token from above)
curl -s "http://127.0.0.1:7350/v2/account" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq .
```

## Step 5: Development Workflow

1. **Edit your Go code** in `modules/main.go`
2. **Rebuild the module**: 
   ```bash
   docker run --rm -v $(pwd)/modules:/app -w /app golang:1.21 sh -c "go build -buildmode=plugin -o main.so main.go"
   cp modules/main.so data/modules/
   ```
3. **Restart Nakama**: `docker compose restart nakama`
4. **Test your changes**: `node test-nakama.js`

## Common Issues and Solutions

### Issue: "Plugin not found"
**Solution**: Make sure you're building with the correct flags:
```bash
go build -buildmode=plugin -o main.so main.go
```

### Issue: "invalid ELF header"
**Solution**: Build the module for Linux architecture:
```bash
docker run --rm -v $(pwd)/modules:/app -w /app golang:1.21 sh -c "go build -buildmode=plugin -o main.so main.go"
```

### Issue: "Database connection failed"
**Solution**: Check your database configuration:
```bash
# Check CockroachDB logs
docker compose logs cockroachdb

# Check Nakama logs
docker compose logs nakama
```

### Issue: "Port already in use"
**Solution**: Stop existing containers:
```bash
docker compose down
docker compose up -d
```

## Service Information

### Current Configuration
- **Database**: CockroachDB v24.1.22 (PostgreSQL-compatible)
- **Nakama Version**: 3.26.0
- **API Port**: 7350
- **Console Port**: 7351
- **Database Port**: 26257
- **CockroachDB Admin UI**: http://localhost:8080

### Access Points
- **Nakama API**: http://localhost:7350
- **Nakama Console**: http://localhost:7351
- **CockroachDB Admin**: http://localhost:8080

## Useful Commands

```bash
# View Nakama logs
docker compose logs nakama

# View CockroachDB logs
docker compose logs cockroachdb

# Check service status
docker ps

# Access CockroachDB CLI
docker exec -it effective-golang-cockroachdb-1 cockroach sql --insecure

# Check Nakama health
curl http://localhost:7350/v2/health

# Stop all services
docker compose down

# Start all services
docker compose up -d
```

## Production Considerations

1. **Use proper authentication** (not custom auth)
2. **Set up SSL/TLS certificates**
3. **Configure proper database credentials**
4. **Set up monitoring and logging**
5. **Use environment variables for configuration**
6. **Implement proper error handling**
7. **Add rate limiting and validation**
8. **Use CockroachDB clustering for high availability**

## Next Steps

1. **Read the learning guide**: `docs/nakama-learning-guide.md`
2. **Try the example game**: `tutorials/examples/nakama_simple_game.go`
3. **Test with the client**: `tutorials/examples/nakama_client_example.js`
4. **Explore the Nakama documentation**: https://heroiclabs.com/docs/

## Support

- **Documentation**: https://heroiclabs.com/docs/
- **GitHub**: https://github.com/heroiclabs/nakama
- **Discord**: https://discord.gg/heroiclabs
- **Forums**: https://forum.heroiclabs.com/
