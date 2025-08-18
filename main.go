package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Effective Go Learning Project")
	fmt.Println("=============================")
	fmt.Println()
	fmt.Println("This project demonstrates Go fundamentals and best practices.")
	fmt.Println()
	fmt.Println("To run the server application:")
	fmt.Println("  go run cmd/server/*.go")
	fmt.Println()
	fmt.Println("To run the tutorial examples:")
	fmt.Println("  go run tutorials/examples/variables.go")
	fmt.Println("  go run tutorials/examples/basic_usage.go")
	fmt.Println()
	fmt.Println("To run tests:")
	fmt.Println("  go test ./...")
	fmt.Println()
	fmt.Println("To run benchmarks:")
	fmt.Println("  go test -bench=. ./...")
	fmt.Println()
	fmt.Println("For more information, see the README.md file.")
	
	// Exit with success
	os.Exit(0)
}