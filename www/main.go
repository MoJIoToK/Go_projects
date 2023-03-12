package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type User struct {
	FirstName             string
	LastName              string
	Age                   uint16
	Money                 int16
	Avg_grades, Happiness float64
	Desription            string
	Hobbies               []string
}

func (u User) getAllInfo() string {
	return fmt.Sprintf("Username is: %s %s. He is %d years old and he has %d money",
		u.FirstName, u.LastName, u.Age, u.Money)
}

func (u *User) setNewName(newName string) {
	u.FirstName = newName
}

func home_page(w http.ResponseWriter, r *http.Request) {
	nick := User{"Nick", "Sabjur", 26, 10, 4.0, 0.8, "developer", []string{"Footbal", "Computer"}}
	// //nick.setNewName("Alex")
	// fmt.Fprintf(w, `<h1>Main Text</h1>
	// <b>Main Text</b>`)
	tmpl, _ := template.ParseFiles("C:/Users/Nick/Desktop/Коля/Go/www/templates/home_page.html")
	tmpl.Execute(w, nick)
}

func contacts_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is contacts page")
}

func handleRequest() {
	http.HandleFunc("/", home_page)   //отслеживает переход по определенному URL адресу и вызывает необходимую функцию
	http.ListenAndServe(":8080", nil) // nil - null(none)
	http.HandleFunc("/contacts/", contacts_page)
}

func main() {
	//var nick User = ...

	handleRequest()
}
