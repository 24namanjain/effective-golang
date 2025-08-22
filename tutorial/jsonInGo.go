package main

import (
	"encoding/json"
	"fmt"
)

func jsonEncodingDecoding() {

	// create a struct
	type Person struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	// create a person
	person := Person{Name: "John", Age: 30, Email: "john@example.com"}

	// encode the person to json
	jsonPerson, err := json.Marshal(person)
	if err != nil {
		fmt.Println("Error encoding person to json: ", err)
	}

	fmt.Println("jsonPerson: ", string(jsonPerson))

	// decode the json to a person
	var decodedPerson Person
	err = json.Unmarshal(jsonPerson, &decodedPerson)
	if err != nil {
		fmt.Println("Error decoding person from json: ", err)
	}

	fmt.Println("decodedPerson: ", decodedPerson)
}
