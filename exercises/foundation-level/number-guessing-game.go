//go:build hello2

package main

import "fmt"

var low = 1
var high = 1000

func confirmUserNum(guess int) int {

	fmt.Printf("Is your number %d? (0 for yes, 1 for too low, 2 for too high)\n", guess)

	var confirm int
	fmt.Scanf("%d", &confirm)

	if confirm != 0 && confirm != 1 && confirm != 2 {
		// throw an exception
		panic("Invalid input")
	}

	return confirm

}

func guessUserNum(tries int) bool {

	for i := low; i < tries; i++ {

		mid := (low + high) / 2

		var success = confirmUserNum(mid)

		if success == 0 {
			return true
		}

		if success == 1 {
			low = mid + 1
		} else {
			high = mid - 1
		}

	}

	return false

}

func main() {

	var tries = 15

	fmt.Printf("Think of a number between %d to %d\n", low, high)

	success := guessUserNum(tries)

	if success {
		fmt.Println("I guessed your number!")
	} else {
		fmt.Println("I couldn't guess your number!")
	}

}
