package main

import (
	"fmt"
	"os"
)

func errorHandlingDemo() {
	// Example 1: Basic error handling with custom error
	result, err := divide(10, 0)
	if err != nil {
		// Error has context from divide()
		fmt.Println("Error occurred:", err)
	} else {
		fmt.Println("Result:", result)
	}

	// Example 2: Error handling with file open (simulated)
	err = doSomethingThatFails()
	if err != nil {
		// Error has context from doSomethingThatFails()
		fmt.Println("Operation failed:", err)
	} else {
		fmt.Println("Operation succeeded")
	}
}

// divide returns the result of a / b, or an error if b == 0
func divide(a, b int) (int, error) {
	if b == 0 {
		// divide(10, 0): invalid argument
		// os.ErrInvalid is a constant error that is returned when an invalid argument is passed to a function
		// %w is a placeholder for the error that is returned
		// Wrap the error with context
		return 0, fmt.Errorf("divide(%d, %d): %w", a, b, os.ErrInvalid)
	}
	return a / b, nil
}

// doSomethingThatFails simulates a function that returns an error
func doSomethingThatFails() error {
	baseErr := fmt.Errorf("something went wrong")
	// Wrap with higher-level context
	return fmt.Errorf("doSomethingThatFails: %w", baseErr)
}
