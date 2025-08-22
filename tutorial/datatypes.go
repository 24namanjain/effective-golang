package main

import "fmt"

func goDatatypes() {
	boolDatatypes()

	intDatatypes()

	floatDatatypes()

	arrayDatatypes()

	structDatatypes()

	mapDatatypes()
}

func boolDatatypes() {

	var var1 bool = false
	var var2 bool = false

	fmt.Println("var1: ", var1)
	fmt.Println("var2: ", var2)

	// else must be on the same line as the closing braces of if statement
	// as go insert semicolon
	if var1 {
		fmt.Println("var1 is true")
	} else {
		fmt.Println("var1 is false")
	}

}

func intDatatypes() {

	// unsigned integers
	var u8 uint8 = 255                    // 8 bits
	var u16 uint16 = 65535                // 16 bits
	var u32 uint32 = 4294967295           // 32 bits
	var u64 uint64 = 18446744073709551615 // 64 bits

	// signed integers
	var i8 int8 = -128                  // 8 bits
	var i16 int16 = 32767               // 16 bits
	var i32 int32 = 2147483647          // 32 bits
	var i64 int64 = 9223372036854775807 // 64 bits

	fmt.Println("u8: ", u8)
	fmt.Println("u16: ", u16)
	fmt.Println("i8: ", i8)
	fmt.Println("i16: ", i16)

	fmt.Println("u32: ", u32)
	fmt.Println("u64: ", u64)
	fmt.Println("i32: ", i32)
	fmt.Println("i64: ", i64)

}

func floatDatatypes() {

	var pi float32 = 3.14159
	var e float64 = 2.718281828

	fmt.Println("pi: ", pi)
	fmt.Println("e: ", e)

	var firstComplexNum complex64 = 1 + 2i
	var secondComplexNum complex128 = 3 + 4i

	fmt.Println("firstComplexNum: ", firstComplexNum)
	fmt.Println("secondComplexNum: ", secondComplexNum)

	fmt.Println("real and img parts of firstComplexNum: ", real(firstComplexNum), imag(firstComplexNum))
	fmt.Println("real and img parts of secondComplexNum: ", real(secondComplexNum), imag(secondComplexNum))

}

func arrayDatatypes() {

	var arr [3]int = [3]int{10, 22, 31}

	fmt.Println(arr)

	// print the length of the array
	fmt.Println("length of arr: ", len(arr))

	// print the first element of the array
	fmt.Println("first element of arr: ", arr[0])

	// print the elements of the array
	fmt.Println("elements of arr: ", arr)

	// print the elements of the array using a for loop with index
	for i := 0; i < len(arr); i++ {
		fmt.Println("element ", i, " of arr: ", arr[i])
	}

	// print the elements of the array using a for loop with index and element
	for _, element := range arr {
		fmt.Println("element of arr: ", element)
	}

}

func structDatatypes() {

	type Person struct {
		name  string
		age   int
		email string
	}

	var person Person = Person{name: "John", age: 30, email: "john@example.com"}

	fmt.Println("person: ", person)

	// print the name of the person
	fmt.Println("name of person: ", person.name)

	// print the age of the person
	fmt.Println("age of person: ", person.age)

	// print the email of the person
	fmt.Println("email of person: ", person.email)

	// print the name, age and email of the person
	fmt.Println("name, age and email of person: ", person.name, person.age, person.email)
}

func mapDatatypes() {

	var map1 map[string]int = map[string]int{"apple": 1, "banana": 2, "cherry": 3}

	fmt.Println("map1: ", map1)

	// print the value of the key "apple"
	fmt.Println("value of key apple: ", map1["apple"])

	fmt.Println("length of map1: ", len(map1))
}
