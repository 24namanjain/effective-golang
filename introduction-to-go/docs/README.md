# Nakama Documentation

This directory contains comprehensive documentation for setting up and using Nakama with your Go applications.

## 📚 Documentation Files

### 1. **nakama-setup-guide.md**
A step-by-step guide to get Nakama running with CockroachDB using Docker Compose.

**What you'll learn:**
- ✅ Install Nakama with CockroachDB
- ✅ Create and build Go modules
- ✅ Test the setup with JavaScript client
- ✅ Common issues and solutions
- ✅ Development workflow

**Quick Start:**
```bash
# Start Nakama
docker compose up -d

# Test the setup
node test-nakama.js
```

### 2. **nakama-learning-guide.md**
Comprehensive learning guide covering all Nakama concepts and features.

**What you'll learn:**
- 🎯 Core concepts (RPC, Events, Leaderboards, Matches)
- 🏗️ Architecture and design patterns
- 💻 Complete code examples
- 🔧 Best practices and common patterns
- 🚀 Real-world implementation examples

### 3. **nakama-quick-reference.md**
Quick reference guide for running and validating Nakama.

**What you'll learn:**
- 🚀 How to run and validate Nakama
- 📋 Common commands and troubleshooting
- 🛠️ Development workflow
- 📊 Service information and status

### 4. **tutorials/** (Organized Learning Materials)
Comprehensive tutorials organized by topic:

#### **📁 tutorials/go-basics/**
- [01-getting-started.md](tutorials/go-basics/01-getting-started.md) - Setup and first steps
- [02-go-syntax-basics.md](tutorials/go-basics/02-go-syntax-basics.md) - Fundamental syntax
- [03-data-structures.md](tutorials/go-basics/03-data-structures.md) - Data types and structures

#### **📁 tutorials/go-advanced/**
- [04-error-handling.md](tutorials/go-advanced/04-error-handling.md) - Error handling patterns
- [06-core-concepts.md](tutorials/go-advanced/06-core-concepts.md) - Core Go concepts
- [07-concurrency.md](tutorials/go-advanced/07-concurrency.md) - Goroutines and channels

#### **📁 tutorials/nakama/**
- [nakama-learning-guide.md](tutorials/nakama/nakama-learning-guide.md) - Complete Nakama guide

#### **📁 tutorials/project/**
- [05-project-overview.md](tutorials/project/05-project-overview.md) - Project architecture
- [golang-req.md](tutorials/project/golang-req.md) - Learning requirements

#### **📁 tutorials/ (Root)**
- [00-index.md](tutorials/00-index.md) - Complete tutorial index
- [SEQUENTIAL_LEARNING.md](tutorials/SEQUENTIAL_LEARNING.md) - Learning path guide

## 🎮 Current Setup Status

Your Nakama server is currently running with:
- **Database**: CockroachDB v24.1.22
- **Nakama Version**: 3.26.0
- **Status**: ✅ Fully functional
- **Access**: http://localhost:7350 (API), http://localhost:7351 (Console)

## 🚀 Quick Commands

```bash
# Start services
docker compose up -d

# View logs
docker compose logs nakama

# Test connection
node test-nakama.js

# Stop services
docker compose down
```

## 📖 Learning Path

### **For Nakama Game Development:**
1. **Start with**: `nakama-setup-guide.md` - Get everything running
2. **Then read**: `nakama-learning-guide.md` - Understand concepts
3. **Use**: `nakama-quick-reference.md` - Quick commands and validation
4. **Practice with**: Examples in `tutorials/examples/`
5. **Build your game**: Using the patterns you've learned

### **For Go Learning:**
1. **Start with**: [tutorials/00-index.md](tutorials/00-index.md) - Complete overview
2. **Follow**: [tutorials/SEQUENTIAL_LEARNING.md](tutorials/SEQUENTIAL_LEARNING.md) - Structured path
3. **Learn basics**: [tutorials/go-basics/](tutorials/go-basics/) - Fundamental concepts
4. **Master advanced**: [tutorials/go-advanced/](tutorials/go-advanced/) - Advanced patterns
5. **Study project**: [tutorials/project/](tutorials/project/) - Architecture and requirements

## 🔗 Useful Links

- **Official Docs**: https://heroiclabs.com/docs/
- **GitHub**: https://github.com/heroiclabs/nakama
- **Discord**: https://discord.gg/heroiclabs
- **Forums**: https://forum.heroiclabs.com/

## 🛠️ Development Workflow

1. **Edit Go modules** in `modules/main.go`
2. **Rebuild**: `docker run --rm -v $(pwd)/modules:/app -w /app golang:1.21 sh -c "go build -buildmode=plugin -o main.so main.go"`
3. **Deploy**: `cp modules/main.so data/modules/ && docker compose restart nakama`
4. **Test**: `node test-nakama.js`

## 📝 Notes

- The setup uses **CockroachDB** (not PostgreSQL) as it's the recommended database for Nakama
- All modules must be built for **Linux** architecture when using Docker
- The JavaScript client examples use the official `@heroiclabs/nakama-js` SDK
- All examples are tested and working with the current configuration
