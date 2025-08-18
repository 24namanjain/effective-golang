# Core Go Concepts

This document explains the fundamental Go concepts demonstrated in this project.

## 1. Structs and Methods

### Structs
Structs are composite data types that group together variables under a single name.

```go
type User struct {
    ID        string    `json:"id" db:"id"`
    Username  string    `json:"username" db:"username"`
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"-" db:"password"` // "-" excludes from JSON
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    IsActive  bool      `json:"is_active" db:"is_active"`
}
```

**Key Points:**
- Use struct tags for JSON marshaling/unmarshaling
- Use `json:"-"` to exclude fields from JSON output
- Group related fields together
- Use meaningful field names

### Methods
Methods are functions with a receiver that operate on a specific type.

```go
func (u *UserStats) GetWinRate() float64 {
    if u.TotalGames == 0 {
        return 0.0
    }
    return float64(u.Wins) / float64(u.TotalGames) * 100
}
```

**Key Points:**
- Use pointer receivers (`*Type`) when you need to modify the struct
- Use value receivers (`Type`) for read-only operations
- Keep methods focused and cohesive

## 2. Interfaces

Interfaces define behavior without implementation details.

```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}
```

**Key Points:**
- Interfaces should be small and focused
- Use interfaces for dependency injection
- Implement interfaces implicitly (no explicit declaration needed)
- Follow the "accept interfaces, return concrete types" principle

## 3. Error Handling

Go uses explicit error handling with return values.

```go
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
    // Check if user already exists
    existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
    if err == nil && existingUser != nil {
        return nil, fmt.Errorf("username already taken: %w", ErrUserAlreadyExists)
    }
    
    // Create new user
    user, err := models.NewUser(req.Username, req.Email, req.Password)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}
```

**Key Points:**
- Always check errors explicitly
- Use `fmt.Errorf("context: %w", err)` for error wrapping
- Create custom error types for domain-specific errors
- Return errors early to avoid deep nesting

## 4. Maps and Slices

### Maps
Maps are key-value data structures.

```go
// In-memory storage example
type InMemoryUserRepository struct {
    users map[string]*User
    mutex sync.RWMutex
}

// Usage
user, exists := r.users[id]
if !exists {
    return nil, ErrUserNotFound
}
```

### Slices
Slices are dynamic arrays.

```go
// Creating slices
entries := make([]LeaderboardEntry, 0, len(leaderboard.Entries))

// Appending to slices
entries = append(entries, newEntry)

// Slicing
result := entries[:count]
```

**Key Points:**
- Use `make()` to pre-allocate slices with capacity
- Use `append()` to add elements
- Be careful with slice operations (they share underlying array)

## 5. JSON Handling

Go provides built-in JSON marshaling/unmarshaling.

```go
// Marshaling (struct to JSON)
data, err := json.Marshal(user)
if err != nil {
    return err
}

// Unmarshaling (JSON to struct)
var user User
err := json.Unmarshal(data, &user)
if err != nil {
    return err
}
```

**Key Points:**
- Use struct tags to control JSON field names
- Use `json:"omitempty"` to exclude zero values
- Use `json:"-"` to exclude fields entirely
- Always check for marshaling/unmarshaling errors

## 6. Context

Context provides request-scoped values, cancellation, and deadlines.

```go
func (s *GameService) CreateGame(ctx context.Context, player1ID, player2ID string) (*Game, error) {
    // Validate players exist
    player1, err := s.userRepo.GetByID(ctx, player1ID)
    if err != nil {
        return nil, fmt.Errorf("player1 not found: %w", err)
    }
    
    // Check for cancellation
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        // Continue processing
    }
    
    return game, nil
}
```

**Key Points:**
- Always pass context as the first parameter
- Check for context cancellation in long-running operations
- Use context for timeouts and deadlines
- Don't store context in structs

## 7. Packages and Imports

### Package Organization
```
internal/
├── auth/          # Authentication functionality
├── game/          # Game logic
├── leaderboard/   # Leaderboard management
└── models/        # Data structures

pkg/
└── utils/         # Public utility functions
```

### Import Guidelines
```go
import (
    // Standard library first
    "context"
    "fmt"
    "time"
    
    // Third-party packages
    "github.com/gorilla/mux"
    
    // Internal packages
    "effective-golang/internal/models"
    "effective-golang/pkg/utils"
)
```

**Key Points:**
- Use `internal/` for private application code
- Use `pkg/` for public libraries
- Group imports by standard library, third-party, and internal
- Use meaningful package names

## 8. Naming Conventions

### Packages
- Use lowercase, no underscores
- Keep names short and descriptive
- Examples: `auth`, `game`, `leaderboard`

### Types and Interfaces
- Use PascalCase
- Be descriptive
- Examples: `UserRepository`, `GameService`, `LeaderboardEntry`

### Variables and Functions
- Use camelCase
- Be descriptive but concise
- Examples: `userID`, `createGame`, `getTopEntries`

### Constants
- Use ALL_CAPS with underscores
- Examples: `MAX_RETRY_COUNT`, `DEFAULT_TIMEOUT`

### Receivers
- Use short, meaningful names
- Examples: `func (u *User)`, `func (s *Service)`

## 9. Best Practices

### Code Organization
- Keep functions small and focused
- Use meaningful variable names
- Add comments for complex logic
- Follow the single responsibility principle

### Error Handling
- Always check errors
- Provide meaningful error messages
- Use custom error types for domain errors
- Don't ignore errors

### Performance
- Use appropriate data structures
- Pre-allocate slices when possible
- Use pointers for large structs
- Profile your code for bottlenecks

### Testing
- Write tests for all public functions
- Use table-driven tests
- Test error conditions
- Use benchmarks for performance-critical code

## 10. Common Patterns

### Repository Pattern
```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
    // ... other methods
}
```

### Service Layer
```go
type AuthService struct {
    userRepo  UserRepository
    cacheRepo CacheRepository
}
```

### Dependency Injection
```go
func NewAuthService(userRepo UserRepository, cacheRepo CacheRepository) *AuthService {
    return &AuthService{
        userRepo:  userRepo,
        cacheRepo: cacheRepo,
    }
}
```

These patterns promote:
- Testability
- Loose coupling
- Code reusability
- Maintainability
