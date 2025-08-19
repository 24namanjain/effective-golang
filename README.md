# Effective Golang

A comprehensive learning project covering Go programming, design patterns, and modern backend technologies.

## Project Structure

```
effective-golang/
├── introduction-to-go/          # Basic Go learning project
│   ├── cmd/                     # Command-line applications
│   ├── internal/                # Internal packages
│   ├── pkg/                     # Public packages
│   ├── docs/                    # Go documentation
│   └── README.md               # Go learning guide
│
├── slack-notifier/              # Complete Slack notification system
│   ├── cmd/server/             # Main application
│   ├── internal/               # Internal modules
│   │   ├── config/            # Configuration management
│   │   ├── events/            # Event types and builders
│   │   └── notifier/          # Core notification service
│   ├── pkg/slack/             # Slack API client
│   ├── examples/              # Usage examples
│   ├── docs/                  # Documentation
│   └── README.md              # Project documentation
│
└── nakama-learning/            # Nakama game server learning
    ├── docs/                  # Nakama documentation
    ├── examples/              # Code examples
    ├── modules/               # Go modules for Nakama
    ├── data/modules/          # Compiled modules
    ├── docker-compose.yml     # Nakama + Database setup
    └── README.md              # Nakama learning guide
```

## What Each Project Teaches

### 1. **Introduction to Go** (`introduction-to-go/`)
**Focus**: Core Go programming concepts and design patterns

**Topics Covered**:
- Go syntax and fundamentals
- Design patterns (Repository, Unit of Work, Strategy, etc.)
- Package management and project structure
- Testing and documentation
- Best practices and idioms

**Key Files**:
- `cmd/` - Command-line applications
- `internal/models/` - Domain models and interfaces
- `pkg/utils/` - Reusable utilities
- `docs/` - Learning materials and examples

### 2. **Slack Notifier** (`slack-notifier/`)
**Focus**: Building a production-ready Go application

**Topics Covered**:
- Modular architecture design
- Configuration management
- Event-driven programming
- HTTP API development
- Concurrent processing with goroutines
- Error handling and logging
- Testing and documentation

**Key Features**:
- Multiple event types (system, user, business, alerts)
- Severity levels (info, warning, error, critical)
- Rich Slack message formatting
- Worker pool for concurrent processing
- HTTP API for triggering notifications
- Comprehensive testing suite

**Architecture**:
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   HTTP API  │    │   Events    │    │   Workers   │
│   Server    │───▶│   Queue     │───▶│   Pool      │
└─────────────┘    └─────────────┘    └─────────────┘
                           │                   │
                           ▼                   ▼
                    ┌─────────────┐    ┌─────────────┐
                    │   Event     │    │   Slack     │
                    │   Builder   │    │   Client    │
                    └─────────────┘    └─────────────┘
```

### 3. **Nakama Learning** (`nakama-learning/`)
**Focus**: Game server development and real-time applications

**Topics Covered**:
- Nakama game server concepts
- Real-time communication (WebSockets)
- RPC (Remote Procedure Call) development
- Leaderboards and scoring systems
- User authentication and sessions
- Database integration (CockroachDB)
- Multiplayer game architecture

**Key Features**:
- Complete multiplayer game example
- Real-time chat system
- Leaderboard functionality
- User data storage
- JavaScript client examples
- Docker-based development environment

**Architecture**:
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │    │   Client    │    │   Client    │
│  (Mobile/   │    │   (Web)     │    │  (Desktop)  │
│   Console)  │    │             │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                    ┌─────────────┐
                    │   Nakama    │
                    │   Server    │
                    │             │
                    │ ┌─────────┐ │
                    │ │Database │ │
                    │ │(Cockroach│ │
                    │ │   DB)   │ │
                    │ └─────────┘ │
                    └─────────────┘
```

## Learning Path

### For Beginners
1. **Start with `introduction-to-go/`** - Learn Go fundamentals
2. **Move to `slack-notifier/`** - Build a real application
3. **Explore `nakama-learning/`** - Learn game server development

### For Intermediate Developers
1. **Focus on `slack-notifier/`** - Study the architecture and patterns
2. **Experiment with `nakama-learning/`** - Build multiplayer features
3. **Review `introduction-to-go/`** - Refine Go skills

### For Advanced Developers
1. **Extend `slack-notifier/`** - Add new event types, integrations
2. **Build games with `nakama-learning/`** - Create complex multiplayer games
3. **Contribute patterns to `introduction-to-go/`** - Share knowledge

## Technology Stack

### Core Technologies
- **Go 1.21+** - Primary programming language
- **Docker** - Containerization and deployment
- **Git** - Version control

### Slack Notifier
- **slack-go/slack** - Slack API client
- **godotenv** - Environment configuration
- **net/http** - HTTP server and client
- **context** - Request cancellation and timeouts
- **goroutines** - Concurrent processing

### Nakama Learning
- **Nakama** - Game server framework
- **CockroachDB** - Distributed database
- **WebSockets** - Real-time communication
- **JavaScript** - Client-side development

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Docker and Docker Compose
- Node.js (for Nakama client examples)
- Git

### Quick Start

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd effective-golang
   ```

2. **Choose a project to start with**
   ```bash
   # For Go basics
   cd introduction-to-go
   
   # For Slack notifications
   cd slack-notifier
   
   # For game server development
   cd nakama-learning
   ```

3. **Follow the project-specific README**
   Each project has its own detailed documentation and setup instructions.

## Contributing

1. **Fork the repository**
2. **Create a feature branch**
3. **Make your changes**
4. **Add tests and documentation**
5. **Submit a pull request**

## Project Goals

- **Learning**: Provide comprehensive learning materials for Go development
- **Practical**: Focus on real-world applications and use cases
- **Modular**: Each project is self-contained and focused
- **Extensible**: Easy to extend and build upon
- **Production-ready**: Demonstrate best practices and patterns

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Note**: Each project is designed to be independent. You can work on any project without needing to understand the others, but they build upon each other in complexity and scope.
