# Project Structure

This project has been restructured to separate the main application from tutorial/learning materials.

## 📁 Directory Structure

```
effective-golang/
├── cmd/                    # Application entry points
│   └── server/            # Main server application
│       ├── main.go        # Server entry point
│       └── handlers.go    # HTTP request handlers
├── internal/              # Private application code
│   ├── auth/              # Authentication functionality
│   │   └── service.go     # Auth service implementation
│   ├── game/              # Game logic
│   │   └── service.go     # Game service implementation
│   ├── leaderboard/       # Leaderboard management
│   │   └── service.go     # Leaderboard service implementation
│   └── models/            # Data structures and interfaces
│       ├── user.go        # User model and validation
│       ├── game.go        # Game model and logic
│       ├── leaderboard.go # Leaderboard model
│       └── repository.go  # Repository interfaces
├── pkg/                   # Public libraries
│   └── utils/             # Utility functions
│       ├── inmemory.go    # In-memory implementations
│       └── response.go    # HTTP response utilities
├── tests/                 # Integration tests
│   └── user_test.go       # User model tests
├── tutorials/             # Learning materials and examples
│   ├── examples/          # Code examples
│   │   ├── variables.go   # Variable declaration examples
│   │   └── basic_usage.go # Basic Go usage patterns
│   ├── docs/              # Documentation
│   │   ├── core-concepts.md    # Core Go concepts
│   │   └── concurrency.md      # Concurrency patterns
│   ├── golang-req.md      # Learning requirements
│   └── README.md          # Tutorial documentation
├── main.go                # Project overview and instructions
├── README.md              # Main project documentation
├── go.mod                 # Go module definition
└── go.sum                 # Go module checksums
```

## 🎯 Separation of Concerns

### Main Application (`cmd/`, `internal/`, `pkg/`, `tests/`)
- **Production-ready code** following Go best practices
- **Clean architecture** with proper separation of concerns
- **Comprehensive testing** with integration tests
- **Real-world patterns** like repository pattern, dependency injection
- **Concurrent processing** with goroutines and channels
- **Error handling** with proper error types and wrapping

### Tutorial Materials (`tutorials/`)
- **Learning examples** for understanding Go concepts
- **Documentation** explaining core concepts and patterns
- **Basic demonstrations** of language features
- **Educational content** separate from production code

## 🚀 Running the Project

### Main Application
```bash
# Run the server
go run cmd/server/*.go

# Run tests
go test ./...

# Build everything
go build ./...
```

### Tutorial Examples
```bash
# Run variable examples
go run tutorials/examples/variables.go

# Run basic usage examples
go run tutorials/examples/basic_usage.go
```

## 📚 Learning Path

1. **Start with tutorials** (`tutorials/`)
   - Read `tutorials/golang-req.md` for learning objectives
   - Practice with examples in `tutorials/examples/`
   - Study documentation in `tutorials/docs/`

2. **Study the main application**
   - Examine the clean architecture in `internal/`
   - Understand patterns in `cmd/server/`
   - Learn from tests in `tests/`

3. **Apply concepts**
   - Use the patterns in your own projects
   - Follow the best practices demonstrated
   - Build upon the foundation provided

## 🎯 Benefits of This Structure

- **Clear separation** between learning and production code
- **Focused learning** with dedicated tutorial materials
- **Professional structure** for the main application
- **Easy navigation** with logical organization
- **Scalable design** that can grow with your needs
- **Best practices** demonstrated throughout

This structure makes it easy to focus on learning Go concepts while also having a reference implementation of production-ready code to study and learn from.
