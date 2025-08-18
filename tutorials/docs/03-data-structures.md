# Data Structures in Go

Data structures are ways to organize and store information in your programs. Think of them like different types of containers - each designed for a specific purpose.

## 1. Arrays and Slices

### Arrays
**What it is**: A fixed-size list of items of the same type
**Think of it like**: A bookshelf with a specific number of slots

```go
// Array with 5 integers
var numbers [5]int = [5]int{1, 2, 3, 4, 5}

// Array with 3 strings
var names [3]string = [3]string{"Alice", "Bob", "Charlie"}

// Shorter way to create arrays
numbers := [5]int{1, 2, 3, 4, 5}
names := [3]string{"Alice", "Bob", "Charlie"}
```

**Important**: Arrays have a fixed size that cannot be changed after creation.

### Slices
**What it is**: A flexible list that can grow or shrink
**Think of it like**: A stretchy rubber band that can hold more or fewer items

```go
// Create an empty slice
var numbers []int

// Create a slice with initial values
numbers := []int{1, 2, 3, 4, 5}

// Create a slice with make (specify length and capacity)
numbers := make([]int, 5)  // Length 5, capacity 5
numbers := make([]int, 0, 10)  // Length 0, capacity 10
```

### Working with Slices

```go
package main

import "fmt"

func main() {
    // Create a slice
    fruits := []string{"apple", "banana", "orange"}
    
    // Add items to a slice
    fruits = append(fruits, "grape")
    fruits = append(fruits, "mango", "kiwi")
    
    // Access items by index
    fmt.Println("First fruit:", fruits[0])  // apple
    fmt.Println("Second fruit:", fruits[1]) // banana
    
    // Get the length
    fmt.Println("Number of fruits:", len(fruits))
    
    // Loop through all items
    for i, fruit := range fruits {
        fmt.Printf("%d: %s\n", i+1, fruit)
    }
    
    // Get a portion of the slice (slicing)
    firstTwo := fruits[0:2]  // [apple, banana]
    lastTwo := fruits[len(fruits)-2:]  // [mango, kiwi]
    
    fmt.Println("First two:", firstTwo)
    fmt.Println("Last two:", lastTwo)
}
```

## 2. Maps

**What it is**: A collection of key-value pairs
**Think of it like**: A dictionary where you look up words to find their meanings

```go
// Create an empty map
var scores map[string]int

// Create a map with initial values
scores := map[string]int{
    "Alice": 95,
    "Bob":   87,
    "Charlie": 92,
}

// Create a map with make
scores := make(map[string]int)
```

### Working with Maps

```go
package main

import "fmt"

func main() {
    // Create a map
    studentScores := map[string]int{
        "Alice":   95,
        "Bob":     87,
        "Charlie": 92,
        "Diana":   88,
    }
    
    // Add a new entry
    studentScores["Eve"] = 91
    
    // Update an existing entry
    studentScores["Bob"] = 89
    
    // Access a value
    aliceScore := studentScores["Alice"]
    fmt.Println("Alice's score:", aliceScore)
    
    // Check if a key exists
    score, exists := studentScores["Frank"]
    if exists {
        fmt.Println("Frank's score:", score)
    } else {
        fmt.Println("Frank not found")
    }
    
    // Delete an entry
    delete(studentScores, "Bob")
    
    // Loop through all entries
    for name, score := range studentScores {
        fmt.Printf("%s: %d\n", name, score)
    }
    
    // Get the number of entries
    fmt.Println("Number of students:", len(studentScores))
}
```

## 3. Structs

**What it is**: A custom data type that groups related information together
**Think of it like**: A form with different fields (name, age, email, etc.)

```go
// Define a struct
type Person struct {
    Name    string
    Age     int
    Email   string
    City    string
}

// Create a struct instance
person := Person{
    Name:  "Alice",
    Age:   25,
    Email: "alice@example.com",
    City:  "New York",
}

// Access struct fields
fmt.Println("Name:", person.Name)
fmt.Println("Age:", person.Age)
```

### Working with Structs

```go
package main

import "fmt"

// Define a struct for a book
type Book struct {
    Title     string
    Author    string
    Year      int
    Pages     int
    Available bool
}

func main() {
    // Create some books
    books := []Book{
        {
            Title:     "The Go Programming Language",
            Author:    "Alan Donovan",
            Year:      2015,
            Pages:     380,
            Available: true,
        },
        {
            Title:     "Effective Go",
            Author:    "Google",
            Year:      2020,
            Pages:     150,
            Available: true,
        },
        {
            Title:     "Learning Go",
            Author:    "Jon Bodner",
            Year:      2021,
            Pages:     400,
            Available: false,
        },
    }
    
    // Print all books
    fmt.Println("All Books:")
    for i, book := range books {
        fmt.Printf("%d. %s by %s (%d)\n", 
            i+1, book.Title, book.Author, book.Year)
    }
    
    // Find available books
    fmt.Println("\nAvailable Books:")
    for _, book := range books {
        if book.Available {
            fmt.Printf("- %s\n", book.Title)
        }
    }
    
    // Calculate average pages
    totalPages := 0
    for _, book := range books {
        totalPages += book.Pages
    }
    averagePages := float64(totalPages) / float64(len(books))
    fmt.Printf("\nAverage pages: %.1f\n", averagePages)
}
```

## 4. Pointers

**What it is**: A reference to the memory location of a variable
**Think of it like**: A street address that tells you where to find a house

```go
// Create a variable
name := "Alice"

// Create a pointer to that variable
namePointer := &name

// Get the value that the pointer points to
value := *namePointer

fmt.Println("Original value:", name)
fmt.Println("Pointer address:", namePointer)
fmt.Println("Value from pointer:", value)
```

### Why Use Pointers?

```go
package main

import "fmt"

// Function that modifies a value
func updateAge(age *int) {
    *age = *age + 1
}

// Function that doesn't modify the original value
func updateAgeCopy(age int) int {
    return age + 1
}

func main() {
    // Using pointers to modify values
    age := 25
    fmt.Println("Original age:", age)
    
    updateAge(&age)
    fmt.Println("Age after update:", age)
    
    // Without pointers (creates a copy)
    newAge := updateAgeCopy(age)
    fmt.Println("New age:", newAge)
    fmt.Println("Original age unchanged:", age)
}
```

## 5. Interfaces

**What it is**: A contract that defines what methods a type must have
**Think of it like**: A job description that lists required skills

```go
// Define an interface
type Animal interface {
    Speak() string
    Move() string
}

// Create types that implement the interface
type Dog struct {
    Name string
}

func (d Dog) Speak() string {
    return "Woof!"
}

func (d Dog) Move() string {
    return "Running on four legs"
}

type Bird struct {
    Name string
}

func (b Bird) Speak() string {
    return "Tweet!"
}

func (b Bird) Move() string {
    return "Flying"
}

// Function that works with any Animal
func describeAnimal(animal Animal) {
    fmt.Printf("Animal says: %s\n", animal.Speak())
    fmt.Printf("Animal moves: %s\n", animal.Move())
}
```

## 6. Practice Examples

### Example 1: Student Grade Book

```go
package main

import "fmt"

type Student struct {
    Name   string
    Grades []int
}

func (s Student) Average() float64 {
    if len(s.Grades) == 0 {
        return 0
    }
    
    total := 0
    for _, grade := range s.Grades {
        total += grade
    }
    
    return float64(total) / float64(len(s.Grades))
}

func (s Student) HighestGrade() int {
    if len(s.Grades) == 0 {
        return 0
    }
    
    highest := s.Grades[0]
    for _, grade := range s.Grades {
        if grade > highest {
            highest = grade
        }
    }
    
    return highest
}

func main() {
    students := []Student{
        {
            Name:   "Alice",
            Grades: []int{85, 92, 78, 96, 88},
        },
        {
            Name:   "Bob",
            Grades: []int{72, 85, 91, 79, 83},
        },
        {
            Name:   "Charlie",
            Grades: []int{95, 89, 92, 87, 94},
        },
    }
    
    // Print each student's information
    for _, student := range students {
        fmt.Printf("%s:\n", student.Name)
        fmt.Printf("  Grades: %v\n", student.Grades)
        fmt.Printf("  Average: %.1f\n", student.Average())
        fmt.Printf("  Highest: %d\n", student.HighestGrade())
        fmt.Println()
    }
    
    // Find the student with the highest average
    bestStudent := students[0]
    for _, student := range students {
        if student.Average() > bestStudent.Average() {
            bestStudent = student
        }
    }
    
    fmt.Printf("Best student: %s (Average: %.1f)\n", 
        bestStudent.Name, bestStudent.Average())
}
```

### Example 2: Simple Inventory System

```go
package main

import "fmt"

type Product struct {
    ID       string
    Name     string
    Price    float64
    Quantity int
}

type Inventory struct {
    products map[string]Product
}

func NewInventory() *Inventory {
    return &Inventory{
        products: make(map[string]Product),
    }
}

func (inv *Inventory) AddProduct(product Product) {
    inv.products[product.ID] = product
}

func (inv *Inventory) GetProduct(id string) (Product, bool) {
    product, exists := inv.products[id]
    return product, exists
}

func (inv *Inventory) UpdateQuantity(id string, quantity int) bool {
    product, exists := inv.products[id]
    if !exists {
        return false
    }
    
    product.Quantity = quantity
    inv.products[id] = product
    return true
}

func (inv *Inventory) ListProducts() {
    fmt.Println("Inventory:")
    fmt.Println("ID\tName\t\tPrice\tQuantity")
    fmt.Println("--\t----\t\t-----\t--------")
    
    for _, product := range inv.products {
        fmt.Printf("%s\t%s\t\t$%.2f\t%d\n",
            product.ID, product.Name, product.Price, product.Quantity)
    }
}

func main() {
    // Create inventory
    inventory := NewInventory()
    
    // Add products
    inventory.AddProduct(Product{
        ID:       "P001",
        Name:     "Laptop",
        Price:    999.99,
        Quantity: 10,
    })
    
    inventory.AddProduct(Product{
        ID:       "P002",
        Name:     "Mouse",
        Price:    29.99,
        Quantity: 50,
    })
    
    inventory.AddProduct(Product{
        ID:       "P003",
        Name:     "Keyboard",
        Price:    79.99,
        Quantity: 25,
    })
    
    // List all products
    inventory.ListProducts()
    
    // Update quantity
    inventory.UpdateQuantity("P001", 8)
    
    fmt.Println("\nAfter updating laptop quantity:")
    inventory.ListProducts()
    
    // Get specific product
    if product, exists := inventory.GetProduct("P002"); exists {
        fmt.Printf("\nProduct P002: %s - $%.2f\n", 
            product.Name, product.Price)
    }
}
```

## Common Mistakes to Avoid

1. **Forgetting to initialize maps**: Use `make()` or literal syntax
2. **Accessing slice/map with wrong index**: Check bounds first
3. **Not checking if map key exists**: Use the two-value assignment
4. **Forgetting to use pointers when needed**: Remember when you want to modify the original value
5. **Not handling empty slices/maps**: Always check length before accessing

## Next Steps

1. **Practice**: Try creating your own data structures
2. **Experiment**: Modify the examples and see what happens
3. **Read more**: Check out `04-error-handling.md` for error management
4. **Build something**: Create a simple program using these concepts

Remember: Data structures are the building blocks of programs. Master these, and you'll be able to build complex applications! 🏗️
