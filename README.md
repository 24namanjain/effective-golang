# Effective Golang Learning Project

## 🎯 Project: Multiplayer Game Leaderboard System

A comprehensive learning project that covers Go fundamentals, Nakama backend development, PostgreSQL, and Redis caching.

## 📚 1-Week Intensive Learning Plan

### **Day 1: Go Fundamentals & Project Setup**
**Morning (3 hours):**
- [ ] **Go Basics**: Variables, types, functions, packages
- [ ] **Project Structure**: Understanding `cmd/`, `internal/`, `pkg/` directories
- [ ] **Modules**: `go.mod`, `go.sum`, dependency management
- [ ] **Data Types**: `struct`, `map`, `slice`, `array`, `interface` with JSON tags

**Afternoon (3 hours):**
- [ ] **Error Handling**: Explicit error checking, `fmt.Errorf` with `%w`
- [ ] **HTTP Server**: Basic server with gorilla/mux
- [ ] **JSON Handling**: Struct tags, marshaling/unmarshaling
- [ ] **Naming Conventions**: Packages (lowercase), Structs (PascalCase), vars (camelCase)

**Evening (2 hours):**
- [ ] **Practice**: Build simple REST API with proper error handling
- [ ] **Homework**: Read Go tour (golang.org/tour) sections 1-10

---

### **Day 2: Concurrency & Testing**
**Morning (3 hours):**
- [ ] **Goroutines**: Lightweight threads, `go` keyword
- [ ] **Channels**: Buffered vs unbuffered, send/receive operations
- [ ] **Select Statements**: Multiple channel handling
- [ ] **Context**: Cancellation, deadlines, timeouts

**Afternoon (3 hours):**
- [ ] **Sync Package**: `sync.Mutex`, `sync.WaitGroup`
- [ ] **Worker Pools**: Pattern for high concurrency
- [ ] **Testing**: `testing` package, table-driven tests
- [ ] **Benchmarks**: `go test -bench`, performance measurement

**Evening (2 hours):**
- [ ] **Practice**: Build concurrent worker pool with proper cleanup
- [ ] **Homework**: Write tests for all functions, measure performance

---

### **Day 3: Database & Caching**
**Morning (3 hours):**
- [ ] **PostgreSQL**: Connection pooling, basic CRUD operations
- [ ] **Schema Design**: Normalization, indexes, relationships
- [ ] **Transactions**: ACID properties, rollback scenarios
- [ ] **Query Optimization**: `EXPLAIN ANALYZE`, parameterized queries

**Afternoon (3 hours):**
- [ ] **Redis**: Data types (String, Hash, List, Set, Sorted Set)
- [ ] **Caching Strategies**: TTL, cache invalidation patterns
- [ ] **Connection Management**: Pool limits, error handling
- [ ] **Use Cases**: Session storage, rate limiting, leaderboards

**Evening (2 hours):**
- [ ] **Practice**: Build caching layer with Redis, optimize database queries
- [ ] **Homework**: Design database schema for game leaderboard system

---

### **Day 4: Nakama Backend Fundamentals**
**Morning (3 hours):**
- [ ] **Nakama Architecture**: Server, client, runtime environment
- [ ] **RPC Functions**: `nk.RegisterRpc()`, authentication, validation
- [ ] **Storage Engine**: Key-value store, JSON objects, collections
- [ ] **Events**: `AfterAuthenticate`, `AfterMatchEnd`, custom events

**Afternoon (3 hours):**
- [ ] **Leaderboards**: `nk.LeaderboardCreate()`, `nk.LeaderboardRecordWrite()`
- [ ] **Real-time Events**: WebSocket connections, presence events
- [ ] **Matchmaking**: Basic matchmaker, authoritative matches
- [ ] **Error Handling**: Nakama-specific error patterns

**Evening (2 hours):**
- [ ] **Practice**: Create RPC functions for user registration and score submission
- [ ] **Homework**: Design leaderboard system with proper ranking logic

---

### **Day 5: Advanced Nakama Features**
**Morning (3 hours):**
- [ ] **Match State Management**: In-memory vs Redis persistence
- [ ] **Pub/Sub**: Event broadcasting, real-time updates
- [ ] **Rate Limiting**: Redis-based throttling
- [ ] **Authentication**: Custom auth hooks, session management

**Afternoon (3 hours):**
- [ ] **Advanced Leaderboards**: Multiple leaderboards, time-based rankings
- [ ] **Match Logic**: Game state, player synchronization
- [ ] **Performance**: Connection pooling, query optimization
- [ ] **Monitoring**: Logging, metrics, debugging

**Evening (2 hours):**
- [ ] **Practice**: Build complete matchmaking system with real-time updates
- [ ] **Homework**: Implement rate limiting and monitoring

---

### **Day 6: Integration & Testing**
**Morning (3 hours):**
- [ ] **System Integration**: Connect all components (Go server + Nakama + DB + Redis)
- [ ] **API Design**: RESTful endpoints, proper HTTP status codes
- [ ] **Middleware**: Logging, CORS, authentication
- [ ] **Configuration**: Environment variables, config management

**Afternoon (3 hours):**
- [ ] **Integration Testing**: End-to-end testing, load testing
- [ ] **Error Scenarios**: Network failures, database outages
- [ ] **Performance Testing**: Benchmarking, profiling
- [ ] **Security**: Input validation, SQL injection prevention

**Evening (2 hours):**
- [ ] **Practice**: Run full system tests, identify bottlenecks
- [ ] **Homework**: Document API endpoints and error codes

---

### **Day 7: Project Completion & Best Practices**
**Morning (3 hours):**
- [ ] **Code Review**: Refactor based on Go best practices
- [ ] **Documentation**: API docs, code comments, README
- [ ] **Deployment**: Docker setup, environment configuration
- [ ] **Monitoring**: Health checks, metrics collection

**Afternoon (3 hours):**
- [ ] **Final Testing**: Complete system validation
- [ ] **Performance Optimization**: Identify and fix bottlenecks
- [ ] **Security Audit**: Review for vulnerabilities
- [ ] **Production Readiness**: Error handling, logging, monitoring

**Evening (2 hours):**
- [ ] **Project Presentation**: Demo the complete system
- [ ] **Learning Review**: What worked, what to improve
- [ ] **Next Steps**: Advanced topics, production deployment

## 🏗️ Project Structure

```
effective-golang/
├── cmd/                    # Application entry points
│   └── server/            # Main server application
├── internal/              # Private application code
│   ├── auth/              # Authentication logic
│   ├── game/              # Game logic
│   ├── leaderboard/       # Leaderboard management
│   ├── storage/           # Database operations
│   └── cache/             # Redis caching
├── pkg/                   # Public libraries
├── nakama/                # Nakama server functions
├── tests/                 # Integration tests
├── docs/                  # Documentation
└── scripts/               # Setup and deployment scripts
```

## 🚀 Getting Started

1. **Install Dependencies**
   ```bash
   go mod tidy
   ```

2. **Setup Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Run the Application**
   ```bash
   go run cmd/server/main.go
   ```

## 🎮 Features

- **User Authentication**: Register, login, and session management
- **Game Logic**: Simple game mechanics with scoring
- **Leaderboards**: Real-time leaderboards with rankings
- **Caching**: Redis-based caching for performance
- **Real-time Updates**: WebSocket connections for live updates
- **Matchmaking**: Basic player matching system

## 📖 Learning Objectives

### Go Fundamentals
- ✅ Project structure and modules
- ✅ Error handling patterns
- ✅ Concurrency with goroutines and channels
- ✅ Testing and benchmarking
- ✅ HTTP server development

### Nakama Backend
- ✅ RPC function development
- ✅ Event handling
- ✅ Leaderboard management
- ✅ Real-time communication
- ✅ Matchmaking

### Database & Caching
- ✅ PostgreSQL integration
- ✅ Redis caching strategies
- ✅ Connection pooling
- ✅ Query optimization
- ✅ Transaction management

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

## 📝 Contributing

This is a learning project. Each phase builds upon the previous one, so follow the learning plan sequentially.

## 📚 Daily Learning Resources

### **Day 1 Resources:**
- **Go Tour**: [golang.org/tour](https://golang.org/tour) (Sections 1-10)
- **Go by Example**: [gobyexample.com](https://gobyexample.com)
- **Effective Go**: [golang.org/doc/effective_go.html](https://golang.org/doc/effective_go.html)
- **Practice**: Build a simple calculator API with proper error handling

### **Day 2 Resources:**
- **Concurrency**: [golang.org/doc/effective_go.html#concurrency](https://golang.org/doc/effective_go.html#concurrency)
- **Testing**: [golang.org/pkg/testing/](https://golang.org/pkg/testing/)
- **Practice**: Create a concurrent web scraper with worker pools

### **Day 3 Resources:**
- **PostgreSQL Go Driver**: [github.com/lib/pq](https://github.com/lib/pq)
- **Redis Go Client**: [github.com/redis/go-redis](https://github.com/redis/go-redis)
- **Practice**: Build a caching layer for a user profile system

### **Day 4 Resources:**
- **Nakama Documentation**: [heroiclabs.com/docs/](https://heroiclabs.com/docs/)
- **Nakama Go Runtime**: [heroiclabs.com/docs/runtime-code-basics/](https://heroiclabs.com/docs/runtime-code-basics/)
- **Practice**: Create RPC functions for user management

### **Day 5 Resources:**
- **Nakama Leaderboards**: [heroiclabs.com/docs/gameplay-leaderboards/](https://heroiclabs.com/docs/gameplay-leaderboards/)
- **Nakama Matchmaking**: [heroiclabs.com/docs/gameplay-matchmaker/](https://heroiclabs.com/docs/gameplay-matchmaker/)
- **Practice**: Build a real-time leaderboard with WebSocket updates

### **Day 6 Resources:**
- **Go HTTP**: [golang.org/pkg/net/http/](https://golang.org/pkg/net/http/)
- **Gorilla Mux**: [github.com/gorilla/mux](https://github.com/gorilla/mux)
- **Practice**: Create a complete REST API with authentication

### **Day 7 Resources:**
- **Go Best Practices**: [github.com/golang/go/wiki/CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments)
- **Docker Go**: [docs.docker.com/language/golang/](https://docs.docker.com/language/golang/)
- **Practice**: Deploy your application with Docker

## 🎯 Daily Practice Exercises

### **Day 1 Exercises:**
1. Create a simple struct for a `User` with JSON tags
2. Build a function that returns an error and handle it properly
3. Create a basic HTTP server that serves JSON responses
4. Write a function that marshals/unmarshals JSON data

### **Day 2 Exercises:**
1. Create a goroutine that processes a list of numbers
2. Build a worker pool that processes jobs from a channel
3. Write table-driven tests for your functions
4. Benchmark your functions and optimize them

### **Day 3 Exercises:**
1. Connect to PostgreSQL and perform CRUD operations
2. Create a Redis cache for frequently accessed data
3. Implement connection pooling for both databases
4. Write optimized queries with proper indexing

### **Day 4 Exercises:**
1. Create Nakama RPC functions for user registration
2. Implement leaderboard creation and score submission
3. Handle Nakama events for user authentication
4. Build basic matchmaking logic

### **Day 5 Exercises:**
1. Implement real-time leaderboard updates
2. Create rate limiting using Redis
3. Build match state management with persistence
4. Add monitoring and logging to your functions

### **Day 6 Exercises:**
1. Integrate all components into a working system
2. Write integration tests for the complete flow
3. Implement proper error handling and recovery
4. Add security measures (input validation, SQL injection prevention)

### **Day 7 Exercises:**
1. Refactor code following Go best practices
2. Create comprehensive documentation
3. Set up Docker deployment
4. Perform security audit and performance optimization

## 🔧 Prerequisites & Setup

### **Required Software:**
- Go 1.21+ installed
- PostgreSQL 14+ installed and running
- Redis 6+ installed and running
- Nakama server (can use Docker)
- Git for version control
- A good IDE (VS Code with Go extension recommended)

### **Environment Setup:**
```bash
# Install Go (if not already installed)
# Download from https://golang.org/dl/

# Install PostgreSQL
# macOS: brew install postgresql
# Ubuntu: sudo apt-get install postgresql

# Install Redis
# macOS: brew install redis
# Ubuntu: sudo apt-get install redis-server

# Install Nakama (using Docker)
docker run -d --name nakama -p 7350:7350 -p 7351:7351 heroiclabs/nakama:latest
```

## 📋 Learning Checklist

### **Go Fundamentals:**
- [ ] Variables, types, and functions
- [ ] Structs, interfaces, and methods
- [ ] Error handling patterns
- [ ] Package management with go mod
- [ ] HTTP server development
- [ ] JSON handling and marshaling

### **Concurrency:**
- [ ] Goroutines and channels
- [ ] Context and cancellation
- [ ] Sync package (Mutex, WaitGroup)
- [ ] Worker pools pattern
- [ ] Select statements
- [ ] Race condition prevention

### **Testing:**
- [ ] Unit testing with testing package
- [ ] Table-driven tests
- [ ] Benchmarking
- [ ] Test coverage
- [ ] Mocking and stubbing

### **Database & Caching:**
- [ ] PostgreSQL connection and queries
- [ ] Connection pooling
- [ ] Transactions and ACID
- [ ] Query optimization
- [ ] Redis data types and operations
- [ ] Caching strategies and TTL

### **Nakama Backend:**
- [ ] RPC function development
- [ ] Event handling and hooks
- [ ] Leaderboard management
- [ ] Real-time communication
- [ ] Matchmaking and matches
- [ ] Storage engine usage

### **Production Readiness:**
- [ ] Error handling and logging
- [ ] Security best practices
- [ ] Performance optimization
- [ ] Monitoring and metrics
- [ ] Deployment and Docker
- [ ] Documentation and API specs
