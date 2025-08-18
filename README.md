# Effective Go Learning Project

A comprehensive Go learning project that demonstrates fundamental concepts, best practices, and real-world patterns through a simple game leaderboard system.

## 🎯 Learning Objectives

This project covers all essential Go concepts:

- **Project Structure & Organization**
- **Data Types & Structures**
- **Error Handling**
- **Concurrency & Goroutines**
- **Channels & Communication**
- **Testing & Benchmarking**
- **Naming Conventions**
- **Best Practices**

## 📁 Project Structure

```
effective-golang/
├── cmd/                    # Application entry points
│   └── server/            # Main server application
├── internal/              # Private application code
│   ├── auth/              # Authentication functionality
│   ├── game/              # Game logic
│   ├── leaderboard/       # Leaderboard management
│   └── models/            # Data structures
├── pkg/                   # Public libraries
│   └── utils/             # Utility functions
├── tests/                 # Integration tests
└── tutorials/             # Learning materials and examples
    ├── examples/          # Code examples
    ├── docs/              # Documentation
    └── golang-req.md      # Learning requirements
```

## 🚀 Quick Start

1. **Run the main application:**
   ```bash
   go run cmd/server/*.go
   ```

2. **Run tests:**
   ```bash
   go test ./...
   ```

3. **Run benchmarks:**
   ```bash
   go test -bench=. ./...
   ```

4. **Run tutorial examples:**
   ```bash
   go run tutorials/examples/variables.go
   go run tutorials/examples/basic_usage.go
   ```

## 📚 Learning Path

### 1. **Core Concepts** (`internal/models/`)
- **Structs**: User, Game, Score data structures
- **Interfaces**: Repository pattern for data access
- **Maps & Slices**: Efficient data storage and retrieval

### 2. **Error Handling** (`internal/`)
- **Explicit error checking**: Always check `err != nil`
- **Error wrapping**: Use `fmt.Errorf("context: %w", err)`
- **Custom error types**: Domain-specific errors

### 3. **Concurrency** (`internal/game/`, `internal/leaderboard/`)
- **Goroutines**: Background processing
- **Channels**: Communication between goroutines
- **Context**: Cancellation and timeouts
- **Mutex**: Protecting shared resources

### 4. **Testing** (`tests/`)
- **Table-driven tests**: Efficient test coverage
- **Benchmarks**: Performance testing
- **Integration tests**: End-to-end testing

### 5. **Best Practices**
- **Naming conventions**: Consistent code style
- **Package organization**: Logical grouping
- **Dependency management**: Clean imports

## 🎮 Application Overview

This project implements a simple game leaderboard system with:

- **User Management**: Registration, authentication
- **Game Logic**: Score calculation, game state
- **Leaderboard**: Real-time rankings, caching
- **Concurrent Processing**: Background score updates

## 📖 Detailed Documentation

- [Tutorial Index](./tutorials/docs/00-index.md) - Complete learning guide
- [Getting Started](./tutorials/docs/01-getting-started.md) - Beginner's guide
- [Go Syntax Basics](./tutorials/docs/02-go-syntax-basics.md) - Fundamental syntax
- [Data Structures](./tutorials/docs/03-data-structures.md) - Organizing data
- [Error Handling](./tutorials/docs/04-error-handling.md) - Managing errors
- [Project Overview](./tutorials/docs/05-project-overview.md) - Understanding the application
- [Core Concepts](./tutorials/docs/06-core-concepts.md) - Advanced topics
- [Concurrency](./tutorials/docs/07-concurrency.md) - Parallel processing
- [Learning Requirements](./tutorials/golang-req.md) - Learning objectives

## 🔧 Development

### Prerequisites
- Go 1.24+
- Git

### Setup
```bash
git clone <repository>
cd effective-golang
go mod tidy
```

### Code Style
- Use `gofmt` for formatting
- Follow Go naming conventions
- Write comprehensive tests
- Document public APIs

## 📝 Key Learning Points

1. **Always handle errors explicitly**
2. **Use interfaces for flexibility**
3. **Prefer composition over inheritance**
4. **Write concurrent code safely**
5. **Test everything thoroughly**
6. **Keep packages focused and cohesive**

## 🎯 Next Steps

After completing this project, explore:
- Web frameworks (Gin, Echo)
- Database integration (GORM, sqlx)
- Microservices patterns
- Cloud deployment (Docker, Kubernetes)
- Advanced concurrency patterns

---

**Happy Learning! 🚀**
