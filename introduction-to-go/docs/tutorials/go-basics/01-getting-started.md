# Getting Started with Go

Welcome to Go programming! This guide is designed for complete beginners who have never written Go code before.

## What is Go?

Go (also known as Golang) is a programming language created by Google. It's designed to be:
- **Simple and easy to learn**
- **Fast and efficient**
- **Great for building web services and applications**
- **Excellent for handling multiple tasks at the same time (concurrency)**

## Why Learn Go?

1. **Easy to Learn**: Go has a simple syntax that's easy to understand
2. **Fast Compilation**: Your code compiles quickly, so you can see results fast
3. **Built-in Concurrency**: Go makes it easy to handle multiple tasks simultaneously
4. **Great for Web Development**: Perfect for building APIs and web services
5. **Strong Community**: Lots of resources and helpful developers

## Before You Start

### What You Need to Know
- Basic computer literacy
- Understanding of what programming is
- No prior programming experience required!

### What You'll Learn
1. How to write your first Go program
2. Understanding variables and data types
3. How to control program flow
4. Working with functions
5. Basic error handling

## Your First Go Program

Let's start with the simplest possible Go program. Open the file `tutorials/examples/variables.go` and look at this code:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### What Each Part Means:

- **`package main`**: Every Go program starts with a package declaration. `main` is special - it tells Go this is the starting point of your program.

- **`import "fmt"`**: This brings in a "package" (a collection of useful functions). `fmt` contains functions for printing text to the screen.

- **`func main()`**: This defines a function (a block of code that does something). `main` is special - it's the first function that runs when your program starts.

- **`fmt.Println("Hello, World!")`**: This line prints text to the screen. `Println` means "print and then go to a new line."

## Running Your First Program

1. Open your terminal/command prompt
2. Navigate to the project directory
3. Run: `go run tutorials/examples/variables.go`

You should see output like:
```
Name:  Alice
Age:  20
City:  New York
```

## Understanding Variables

Variables are like containers that hold information. In Go, you can create variables in several ways:

### Method 1: Explicit Declaration
```go
var name string = "Alice"
```
This says: "Create a variable called `name` that holds text (string), and put 'Alice' in it."

### Method 2: Implicit Declaration
```go
var age = 20
```
Go figures out that `age` should be a number (int) because you put a number in it.

### Method 3: Short Declaration
```go
city := "New York"
```
This is a shorter way to create a variable. The `:=` means "create and assign."

## Data Types in Go

Go has several basic data types:

- **`string`**: Text (like "Hello", "Alice", "123")
- **`int`**: Whole numbers (like 1, 42, -7)
- **`float64`**: Decimal numbers (like 3.14, 2.5)
- **`bool`**: True or false values
- **`[]string`**: A list of text items
- **`[]int`**: A list of numbers

## Next Steps

1. **Read the variables example**: Look at `tutorials/examples/variables.go` and try to understand each line
2. **Try modifying the code**: Change the names, ages, and cities to your own values
3. **Run the program**: See how your changes affect the output
4. **Read the documentation**: Check out `core-concepts.md` for more detailed explanations

## Common Beginner Mistakes

1. **Forgetting semicolons**: Go doesn't require semicolons, but you might add them out of habit
2. **Wrong package name**: Always use `package main` for programs you want to run
3. **Missing imports**: If you use `fmt.Println`, you need `import "fmt"`
4. **Case sensitivity**: `Name` and `name` are different in Go

## Getting Help

- **Read the error messages**: Go gives helpful error messages when something goes wrong
- **Check the documentation**: Each concept is explained in detail in the docs folder
- **Experiment**: Try changing things and see what happens
- **Don't worry about understanding everything at once**: Programming is learned step by step

## Practice Exercise

Try this:
1. Open `tutorials/examples/variables.go`
2. Change "Alice" to your name
3. Change the age to your age
4. Change "New York" to your city
5. Run the program and see your information displayed

Congratulations! You've just written and modified your first Go program! 🎉

---

**Next**: Read `02-go-syntax-basics.md` to learn more about Go fundamentals.
