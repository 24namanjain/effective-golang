package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")

	// = is used to assign a value to an already declared variable
	// := is used to declare and assign a value in one step (short variable declaration)

	// Explicit declaration with type and assignment using =
	var a int = 5

	// Short declaration with type inferred as int
	b := 10

	// Swap values of a and b using multiple assignment
	a, b = b, a
	fmt.Println("a =", a, "& b =", b) // a = 10, b = 5

	// Swap again (restores original values)
	a, b = b, a
	fmt.Println("a =", a, "& b =", b) // a = 5, b = 10

	// Short declaration with re-declaration:
	// - b already exists, so it's reassigned
	// - c is new, so it's declared
	b, c := 10, 20
	fmt.Println("b =", b, "& c =", c) // b = 10, c = 20

	// Swap values of b and c
	b, c = c, b
	fmt.Println("b =", b, "& c =", c) // b = 20, c = 10

	// to include and run multiple files
	// use: go run hello.go datatypes.go
	// or use: go run . - this will run all the files in the current directory
	goDatatypes()

	jsonEncodingDecoding()

	errorHandlingDemo()

}
