//go:build hello1

package main

import "fmt"

const PI = 3.14

func main() {

	var helloWorld = "Hello! world"

	fmt.Println(helloWorld)

	fmt.Printf("%s\n", helloWorld)

	var fmtHelloWorld = fmt.Sprintf("Sprintf using %s\n", helloWorld)

	fmt.Println(fmtHelloWorld)

	var message = "This message is from golang"

	fmt.Println(helloWorld + " " + message)

	fmt.Printf(fmt.Sprintf("%s %s", helloWorld, message))

	var name string

	name = getName()
	sayHello(name)

	name = getName()
	sayHello(name)

	age := 10
	address := "India"

	fmt.Printf("%s is %d years old and lives in %s\n", name, age, address)

	fmt.Printf("Pi is %f", PI)

}

func getName() string {
	fmt.Print("\nEnter your name: ")
	var name string
	fmt.Scanln(&name)
	return name
}

func sayHello(name string) {
	fmt.Println("Hello", name)
}
