package main

import "fmt"

func errorHandlingDemo() {
	// Example 1: Basic error handling with custom error
	result, err := divide(10, 0)
	if err != nil {
		fmt.Println("Error occurred:", err)
	} else {
		fmt.Println("Result:", result)
	}

	// Example 2: Error handling with file open (simulated)
	err = doSomethingThatFails()
	if err != nil {
		fmt.Println("Operation failed:", err)
	} else {
		fmt.Println("Operation succeeded")
	}
}

// divide returns the result of a / b, or an error if b == 0
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil
}

// doSomethingThatFails simulates a function that returns an error
func doSomethingThatFails() error {
	return fmt.Errorf("something went wrong")
}
