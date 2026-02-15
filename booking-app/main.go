// go mod init booking-app : this is the command to initialize a new module
package main

import "fmt"

func main() {
	fmt.Println("Hello World")

	var conferenceName = "Go Conference"
	const conferenceTickets = 50
	var remainingTickets uint = 50

	fmt.Printf("conferenceTickets is %T, remainingTickets is %T, conferenceName is %T\n", conferenceTickets, remainingTickets, conferenceName)

	fmt.Printf("Welcome to %v booking application\n", conferenceName)
	fmt.Printf("We have total of %v tickets and %v are still available", conferenceTickets, remainingTickets)
	fmt.Println("Get your tickets here to attend")

	var firstName string
	var lastName string
	var email string
	var userTickets uint

	fmt.Print("Enter your first name: ")
	fmt.Scan(&firstName)

	fmt.Print("Enter your last name: ")
	fmt.Scan(&lastName)

	fmt.Print("Enter your email: ")
	fmt.Scan(&email)

	fmt.Print("Enter number of tickets: ")
	fmt.Scan(&userTickets)

	remainingTickets -= uint(userTickets)

	fmt.Printf("Thank you %v %v for booking %v tickets, you'll recieve a confirmation email @ %v\n", firstName, lastName, email, userTickets)
	fmt.Printf("We've total of %v tickets and %v are remaining", userTickets, remainingTickets)
}
