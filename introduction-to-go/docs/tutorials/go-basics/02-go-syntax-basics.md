# Go Syntax Basics

This guide explains the fundamental syntax rules of Go in simple terms. Think of syntax as the "grammar" of programming - the rules you must follow for Go to understand your code.

## Basic Structure of a Go Program

Every Go program follows this basic pattern:

```go
package main

import "fmt"

func main() {
    // Your code goes here
    fmt.Println("Hello, World!")
}
```

## 1. Package Declaration

**What it is**: Every Go file starts with `package main`
**Why it's needed**: It tells Go what this file is for
**Think of it like**: The title of a book chapter

```go
package main  // This file is a main program
```

## 2. Import Statements

**What it is**: `import` brings in code from other packages
**Why it's needed**: You can't use functions like `fmt.Println` without importing `fmt`
**Think of it like**: Including a library book in your research

```go
import "fmt"           // Import one package
import "strings"       // Import another package

// Or import multiple packages at once:
import (
    "fmt"
    "strings"
    "time"
)
```

## 3. Functions

**What it is**: A block of code that does a specific task
**Why it's needed**: Functions organize your code into reusable pieces
**Think of it like**: A recipe that you can follow multiple times

```go
func main() {
    // This is the main function - it runs first
    fmt.Println("Hello!")
}

func sayHello(name string) {
    // This function takes a name and prints a greeting
    fmt.Println("Hello,", name)
}
```

## 4. Variables and Data Types

### Declaring Variables

Variables are like labeled boxes that store information.

```go
// Method 1: Explicit type declaration
var name string = "Alice"
var age int = 25

// Method 2: Type inference (Go figures out the type)
var city = "New York"
var temperature = 72.5

// Method 3: Short declaration (most common)
name := "Alice"
age := 25
city := "New York"
```

### Basic Data Types

| Type | Example | Description |
|------|---------|-------------|
| `string` | `"Hello"` | Text |
| `int` | `42` | Whole numbers |
| `float64` | `3.14` | Decimal numbers |
| `bool` | `true` | True or false |
| `[]string` | `["a", "b", "c"]` | List of strings |
| `[]int` | `[1, 2, 3]` | List of numbers |

## 5. Control Structures

### If Statements

**What it is**: Code that runs only if a condition is true
**Think of it like**: "If it's raining, take an umbrella"

```go
age := 18

if age >= 18 {
    fmt.Println("You are an adult")
} else {
    fmt.Println("You are a minor")
}

// You can also use else if
if age < 13 {
    fmt.Println("You are a child")
} else if age < 20 {
    fmt.Println("You are a teenager")
} else {
    fmt.Println("You are an adult")
}
```

### For Loops

**What it is**: Code that repeats multiple times
**Think of it like**: "Do this task 10 times"

```go
// Loop 5 times
for i := 0; i < 5; i++ {
    fmt.Println("Count:", i)
}

// Loop through a list
names := []string{"Alice", "Bob", "Charlie"}
for _, name := range names {
    fmt.Println("Hello,", name)
}
```

## 6. Functions with Parameters and Return Values

### Simple Function
```go
func greet(name string) {
    fmt.Println("Hello,", name)
}
```

### Function with Return Value
```go
func add(a int, b int) int {
    return a + b
}

// Usage:
result := add(5, 3)  // result will be 8
```

### Function with Multiple Return Values
```go
func divide(a int, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("cannot divide by zero")
    }
    return a / b, nil
}
```

## 7. Error Handling

**What it is**: Go's way of dealing with things that can go wrong
**Why it's important**: Programs need to handle errors gracefully

```go
result, err := divide(10, 2)
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Result:", result)
```

## 8. Comments

**What it is**: Notes in your code that Go ignores
**Why it's useful**: Explains what your code does

```go
// This is a single-line comment

/*
This is a multi-line comment
It can span multiple lines
*/

func calculateTotal(price float64, tax float64) float64 {
    // Calculate the total including tax
    total := price + (price * tax)
    return total
}
```

## 9. Common Syntax Rules

### Naming Conventions
- **Variables and functions**: Use camelCase (`userName`, `calculateTotal`)
- **Constants**: Use UPPER_CASE (`MAX_SIZE`, `PI`)
- **Packages**: Use lowercase (`fmt`, `strings`)

### Semicolons
- Go automatically adds semicolons, so you don't need to type them
- Don't put semicolons at the end of lines in Go

### Braces
- Opening braces `{` must be on the same line as the statement
- This is required in Go (unlike some other languages)

```go
// Correct:
if x > 0 {
    fmt.Println("Positive")
}

// Wrong:
if x > 0
{
    fmt.Println("Positive")
}
```

## 10. Practice Examples

### Example 1: Simple Calculator
```go
package main

import "fmt"

func main() {
    a := 10
    b := 5
    
    fmt.Println("Addition:", add(a, b))
    fmt.Println("Subtraction:", subtract(a, b))
    fmt.Println("Multiplication:", multiply(a, b))
    
    result, err := divide(a, b)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Division:", result)
    }
}

func add(a, b int) int {
    return a + b
}

func subtract(a, b int) int {
    return a - b
}

func multiply(a, b int) int {
    return a * b
}

func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("cannot divide by zero")
    }
    return a / b, nil
}
```

### Example 2: Working with Lists
```go
package main

import "fmt"

func main() {
    // Create a list of names
    names := []string{"Alice", "Bob", "Charlie", "Diana"}
    
    // Print each name
    fmt.Println("All names:")
    for i, name := range names {
        fmt.Printf("%d: %s\n", i+1, name)
    }
    
    // Find names that start with 'A'
    fmt.Println("\nNames starting with 'A':")
    for _, name := range names {
        if name[0] == 'A' {
            fmt.Println(name)
        }
    }
}
```

## Common Mistakes to Avoid

1. **Forgetting to import packages**: If you use `fmt.Println`, you need `import "fmt"`
2. **Wrong brace placement**: Opening braces must be on the same line
3. **Not handling errors**: Always check `err != nil`
4. **Using wrong variable names**: Go is case-sensitive
5. **Forgetting to use variables**: If you declare a variable, you should use it

## Next Steps

1. **Practice**: Try writing simple programs using these concepts
2. **Read the examples**: Look at `tutorials/examples/basic_usage.go`
3. **Experiment**: Change values and see what happens
4. **Read more**: Check out `03-data-structures.md` for deeper explanations

Remember: Programming is learned by doing! Don't just read - try writing code yourself! 🚀
