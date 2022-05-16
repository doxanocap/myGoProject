package main

import (
	"fmt"
	"net/http"
)

type NewUser struct {
	name                    string
	age                     uint32
	money                   int32
	avgGrade, happinesLevel float64
}

func home_page(w http.ResponseWriter, r *http.Request) {

	user1 := NewUser{"Adam", 25, -50, 4.2, 0.8}

	fmt.Fprintf(w, user1.recieveAllInfo())
}

func (user NewUser) recieveAllInfo() string {
	return fmt.Sprintf("User name is: %s.\n He is %d years old.\n Money left in his debit card %d.\n His GPA %.2f and subjective level of happines %.2f\n", user.name, user.age, user.money, user.avgGrade, user.happinesLevel)
}

func contacts_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is contacs page!")
}

func handleRequest() {
	http.HandleFunc("/", home_page)
	http.HandleFunc("/contacts/", contacts_page)
	http.ListenAndServe(":2004", nil)
}

func main() {

	handleRequest()

}
