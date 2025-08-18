package main;

import "fmt";

func main() {

	// Method 1: explicit declaration
	var name string = "Alice"
	fmt.Println("Name: ", name);

	// Method 2: implicit decalartion
	var age = 20;
	fmt.Println("Age: ", age);


	// Method 3: short declaration
	city := "New York";
	fmt.Println("City: ", city);


	city = "Mumbai";

	// Method 4: multiple variables declaration
	var (
		name2 string = "Bob"
		age2 int = 25
		city2 string = "Los Angeles"
	)

	fmt.Println("Name: ", name, "Age: ", age, "City: ", city);

	fmt.Println("Name2: ", name2, "Age2: ", age2, "City2: ", city2);

	fmt.Println("City after change: ", city); // print the updated city
	
}
