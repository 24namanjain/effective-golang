package main

import "fmt"

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
