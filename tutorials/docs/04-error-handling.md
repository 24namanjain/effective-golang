# Error Handling in Go

Error handling is how your program deals with things that can go wrong. In Go, errors are treated as values, not exceptions. This makes error handling explicit and clear.

## What Are Errors?

**Think of errors like this**: When you're cooking and something goes wrong (burned the food, missing ingredients), you need to handle that situation. Errors in programming are the same - they tell you when something unexpected happened.

## Basic Error Handling

### The Error Type

In Go, errors are just values that implement the `error` interface:

```go
type error interface {
    Error() string
}
```

This means any type that has an `Error()` method that returns a string is an error.

### Creating Simple Errors

```go
package main

import (
    "errors"
    "fmt"
)

func main() {
    // Create a simple error
    err := errors.New("something went wrong")
    fmt.Println("Error:", err)
    
    // Create an error with formatting
    name := "Alice"
    err = fmt.Errorf("user %s not found", name)
    fmt.Println("Error:", err)
}
```

## Common Error Patterns

### 1. Functions That Return Errors

```go
package main

import (
    "errors"
    "fmt"
)

// Function that can return an error
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}

func main() {
    // Call the function and handle the error
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
    
    // Try dividing by zero
    result, err = divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}
```

### 2. Checking for Specific Errors

```go
package main

import (
    "errors"
    "fmt"
)

var (
    ErrNotFound = errors.New("item not found")
    ErrInvalid  = errors.New("invalid input")
)

func findUser(id int) (string, error) {
    if id < 0 {
        return "", ErrInvalid
    }
    
    if id > 100 {
        return "", ErrNotFound
    }
    
    return fmt.Sprintf("User%d", id), nil
}

func main() {
    // Try to find different users
    users := []int{-1, 50, 150}
    
    for _, id := range users {
        name, err := findUser(id)
        if err != nil {
            if errors.Is(err, ErrNotFound) {
                fmt.Printf("User %d not found\n", id)
            } else if errors.Is(err, ErrInvalid) {
                fmt.Printf("Invalid user ID: %d\n", id)
            } else {
                fmt.Printf("Unknown error for user %d: %v\n", id, err)
            }
            continue
        }
        fmt.Printf("Found user: %s\n", name)
    }
}
```

## Error Wrapping

Error wrapping lets you add context to errors while preserving the original error.

```go
package main

import (
    "errors"
    "fmt"
)

func processUser(id int) error {
    // Try to find the user
    name, err := findUser(id)
    if err != nil {
        // Wrap the error with additional context
        return fmt.Errorf("failed to process user %d: %w", id, err)
    }
    
    // Try to save the user
    err = saveUser(name)
    if err != nil {
        return fmt.Errorf("failed to save user %s: %w", name, err)
    }
    
    return nil
}

func findUser(id int) (string, error) {
    if id < 0 {
        return "", errors.New("invalid user ID")
    }
    return fmt.Sprintf("User%d", id), nil
}

func saveUser(name string) error {
    // Simulate a save error
    return errors.New("database connection failed")
}

func main() {
    err := processUser(-1)
    if err != nil {
        fmt.Println("Error:", err)
        
        // Check if the original error was "invalid user ID"
        if errors.Is(err, errors.New("invalid user ID")) {
            fmt.Println("The problem was an invalid user ID")
        }
    }
}
```

## Practical Examples

### Example 1: File Operations

```go
package main

import (
    "fmt"
    "os"
)

func readFile(filename string) (string, error) {
    // Try to open the file
    file, err := os.Open(filename)
    if err != nil {
        return "", fmt.Errorf("failed to open file %s: %w", filename, err)
    }
    defer file.Close() // Make sure file is closed when function ends
    
    // Read the file content (simplified)
    content := "file content here"
    return content, nil
}

func main() {
    content, err := readFile("nonexistent.txt")
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }
    
    fmt.Println("File content:", content)
}
```

### Example 2: User Input Validation

```go
package main

import (
    "errors"
    "fmt"
    "strconv"
)

var (
    ErrEmptyInput = errors.New("input cannot be empty")
    ErrNotNumber  = errors.New("input must be a number")
    ErrTooSmall   = errors.New("number must be at least 1")
    ErrTooLarge   = errors.New("number must be at most 100")
)

func validateAge(input string) (int, error) {
    // Check if input is empty
    if input == "" {
        return 0, ErrEmptyInput
    }
    
    // Try to convert to number
    age, err := strconv.Atoi(input)
    if err != nil {
        return 0, ErrNotNumber
    }
    
    // Check range
    if age < 1 {
        return 0, ErrTooSmall
    }
    
    if age > 100 {
        return 0, ErrTooLarge
    }
    
    return age, nil
}

func main() {
    testInputs := []string{"", "abc", "0", "25", "150"}
    
    for _, input := range testInputs {
        age, err := validateAge(input)
        if err != nil {
            switch {
            case errors.Is(err, ErrEmptyInput):
                fmt.Printf("Input '%s': Please provide an age\n", input)
            case errors.Is(err, ErrNotNumber):
                fmt.Printf("Input '%s': Age must be a number\n", input)
            case errors.Is(err, ErrTooSmall):
                fmt.Printf("Input '%s': Age must be at least 1\n", input)
            case errors.Is(err, ErrTooLarge):
                fmt.Printf("Input '%s': Age must be at most 100\n", input)
            default:
                fmt.Printf("Input '%s': Unknown error: %v\n", input, err)
            }
            continue
        }
        
        fmt.Printf("Input '%s': Valid age %d\n", input, age)
    }
}
```

### Example 3: Database Operations

```go
package main

import (
    "errors"
    "fmt"
)

var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidEmail = errors.New("invalid email format")
    ErrDuplicateUser = errors.New("user already exists")
)

type User struct {
    ID    int
    Name  string
    Email string
}

type Database struct {
    users map[int]User
}

func NewDatabase() *Database {
    return &Database{
        users: make(map[int]User),
    }
}

func (db *Database) CreateUser(user User) error {
    // Check if user already exists
    if _, exists := db.users[user.ID]; exists {
        return ErrDuplicateUser
    }
    
    // Validate email (simplified)
    if user.Email == "" {
        return ErrInvalidEmail
    }
    
    // Save user
    db.users[user.ID] = user
    return nil
}

func (db *Database) GetUser(id int) (User, error) {
    user, exists := db.users[id]
    if !exists {
        return User{}, ErrUserNotFound
    }
    return user, nil
}

func main() {
    db := NewDatabase()
    
    // Try to create users
    users := []User{
        {ID: 1, Name: "Alice", Email: "alice@example.com"},
        {ID: 1, Name: "Bob", Email: "bob@example.com"}, // Duplicate ID
        {ID: 2, Name: "Charlie", Email: ""}, // Invalid email
    }
    
    for _, user := range users {
        err := db.CreateUser(user)
        if err != nil {
            switch {
            case errors.Is(err, ErrDuplicateUser):
                fmt.Printf("Failed to create user %s: User already exists\n", user.Name)
            case errors.Is(err, ErrInvalidEmail):
                fmt.Printf("Failed to create user %s: Invalid email\n", user.Name)
            default:
                fmt.Printf("Failed to create user %s: %v\n", user.Name, err)
            }
            continue
        }
        
        fmt.Printf("Successfully created user: %s\n", user.Name)
    }
    
    // Try to get users
    userIDs := []int{1, 3, 2}
    
    for _, id := range userIDs {
        user, err := db.GetUser(id)
        if err != nil {
            if errors.Is(err, ErrUserNotFound) {
                fmt.Printf("User with ID %d not found\n", id)
            } else {
                fmt.Printf("Error getting user %d: %v\n", id, err)
            }
            continue
        }
        
        fmt.Printf("Found user: %s (%s)\n", user.Name, user.Email)
    }
}
```

## Best Practices

### 1. Always Check Errors

```go
// Good
result, err := someFunction()
if err != nil {
    return err
}

// Bad - ignoring errors
result, _ := someFunction()
```

### 2. Return Errors Early

```go
// Good
func processData(data string) error {
    if data == "" {
        return errors.New("data cannot be empty")
    }
    
    if len(data) < 5 {
        return errors.New("data too short")
    }
    
    // Process the data...
    return nil
}

// Bad - nested conditions
func processData(data string) error {
    if data != "" {
        if len(data) >= 5 {
            // Process the data...
            return nil
        } else {
            return errors.New("data too short")
        }
    } else {
        return errors.New("data cannot be empty")
    }
}
```

### 3. Use Custom Error Types

```go
type ValidationError struct {
    Field string
    Value string
    Rule  string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s (value: %s)", 
        e.Field, e.Rule, e.Value)
}

func validateUser(user User) error {
    if user.Name == "" {
        return ValidationError{
            Field: "name",
            Value: user.Name,
            Rule:  "cannot be empty",
        }
    }
    
    if len(user.Name) < 2 {
        return ValidationError{
            Field: "name",
            Value: user.Name,
            Rule:  "must be at least 2 characters",
        }
    }
    
    return nil
}
```

### 4. Don't Panic (Usually)

```go
// Good - return error
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Bad - panic
func divide(a, b int) int {
    if b == 0 {
        panic("division by zero")
    }
    return a / b
}
```

## Common Mistakes to Avoid

1. **Ignoring errors**: Always check `err != nil`
2. **Not providing context**: Use `fmt.Errorf` to add helpful information
3. **Panicking unnecessarily**: Use panics only for unrecoverable errors
4. **Returning nil errors**: Use `nil` to indicate success, not an empty error
5. **Not handling specific errors**: Use `errors.Is` to check for specific error types

## Next Steps

1. **Practice**: Try writing functions that return errors
2. **Experiment**: Create your own error types
3. **Read more**: Check out `05-project-overview.md` to understand the main application
4. **Build something**: Create a program that handles various error scenarios

Remember: Good error handling makes your programs more reliable and easier to debug! 🛡️
