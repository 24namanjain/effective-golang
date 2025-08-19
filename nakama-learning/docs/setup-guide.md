# Nakama Setup Guide

A step-by-step guide to get Nakama running on your machine.

## Prerequisites

- **Docker** and **Docker Compose** installed
- **Node.js** (for testing)
- **Go** (for building modules)

## Quick Setup

### 1. Start Nakama Server

```bash
# Navigate to the nakama-learning directory
cd nakama-learning

# Start Nakama and database
docker-compose up -d
```

This will start:
- **Nakama server** on port 7352 (HTTP) and 7353 (gRPC)
- **CockroachDB** on port 26258

### 2. Verify It's Running

```bash
# Check if containers are running
docker-compose ps

# View Nakama logs
docker-compose logs -f nakama
```

You should see logs like:
```
{"level":"info","ts":"2024-12-15T14:30:00.000Z","msg":"Node","name":"nakama","version":"3.17.0"}
{"level":"info","ts":"2024-12-15T14:30:01.000Z","msg":"Database connection","addr":"cockroachdb:26257"}
```

### 3. Test Connection

```bash
# Install Node.js dependencies
npm install

# Run the test script
node test-client.js
```

You should see output like:
```
✅ Connected to Nakama
✅ Authentication successful
✅ RPC call successful
```

## What's Running

### Ports
- **7352** - Nakama HTTP API
- **7353** - Nakama gRPC API  
- **26258** - CockroachDB (database)

### Services
- **nakama** - The game server
- **cockroachdb** - The database

## Common Issues

### 1. Port Already in Use
```bash
# Check what's using the ports
lsof -i :7350
lsof -i :26257

# Stop existing services or change ports in docker-compose.yml
```

### 2. Docker Not Running
```bash
# Start Docker Desktop
# Then run:
docker-compose up -d
```

### 3. Permission Issues
```bash
# On Linux/Mac, you might need:
sudo docker-compose up -d
```

## Development Workflow

### 1. **Start Services**
```bash
docker-compose up -d
```

### 2. **Build Go Modules** (if you have custom code)
```bash
# Build the Go module for Linux (Nakama runs on Linux)
docker run --rm -v $(pwd)/modules:/workspace -w /workspace golang:1.21 bash -c "go mod tidy && go build -buildmode=plugin -o main.so main.go"

# Copy to Nakama's module directory
cp modules/main.so data/modules/
```

### 3. **Test Changes**
```bash
node test-client.js
```

### 4. **View Logs**
```bash
docker-compose logs -f nakama
```

### 5. **Stop Services**
```bash
docker-compose down
```

## Configuration

### Environment Variables

You can customize Nakama by setting environment variables in `docker-compose.yml`:

```yaml
environment:
  - NAKAMA_DATABASE_URL=root@cockroachdb:26257/nakama
  - NAKAMA_RUNTIME_PATH=/nakama/data/modules
  - NAKAMA_LOG_LEVEL=info
```

### Database

The default setup uses CockroachDB, but you can switch to PostgreSQL:

```yaml
environment:
  - NAKAMA_DATABASE_URL=postgres://user:password@postgres:5432/nakama
```

## Next Steps

1. **Study the examples** in the `examples/` directory
2. **Modify the test script** to try different features
3. **Build your own modules** in the `modules/` directory
4. **Read the documentation** in the `docs/` directory

## Useful Commands

```bash
# Restart Nakama
docker compose restart nakama

# Reset everything (delete all data)
docker compose down -v
docker compose up -d

# Access Nakama console
curl http://localhost:7352/v2/console

# Check database
docker compose exec cockroachdb cockroach sql --insecure
```

---

**Note**: This is a development setup. For production, you'd need proper security, backups, and scaling configurations.
