package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")

	// to include and run multiple files
	// use: go run hello.go datatypes.go
	// or use: go run . - this will run all the files in the current directory
	goDatatypes()

	jsonEncodingDecoding()

}
