package main

import (
	"fmt"

	"example.com/myapp/internal/auth"

	"github.com/google/uuid"
)

func main() {
	msg := auth.SayHello("Naman")

	fmt.Println(msg)

	id := uuid.New()

	fmt.Println("Generated UUID: ", id)

	datatypes()
	arrays()
	maps()

	datatype_struct()
}

func datatype_struct() {

	type User struct {
		Name string
		Age  int
	}

	var u User = User{Name: "Naman", Age: 28}
	fmt.Println(u)

}

func maps() {

	scores := map[string]int{
		"alice": 90,
		"bob":   88,
	}

	scores["charlie"] = 92

	fmt.Println(scores)
	fmt.Println(scores["alice"])

	val, ok := scores["david"]
	fmt.Println(val, ok)

}

func arrays() {

	// fixed size array
	nums := [3]int{10, 52}

	fmt.Println("nums: ", nums)

	// dynamic array
	fruits := []string{"apple", "banana"}
	fruits = append(fruits, "mango")

	fmt.Println(fruits)
	fmt.Println(fruits[:2]) // slicing

}

func datatypes() {

	var age int = 25

	price := 16.64

	isActive := true

	name := "Naman"

	fmt.Println(name, isActive, price, age)

}
