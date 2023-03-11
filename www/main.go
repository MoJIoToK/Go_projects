package main

import (
	"fmt"
	"net/http"
)

func home_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Go is super easy! ... Maybe")
}

func contacts_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is contacts page")
}

func handleRequest() {
	http.HandleFunc("/", home_page) //отслеживает переход по определенному URL адресу и вызывает необходимую функцию
	http.HandleFunc("/contacts/", contacts_page)
	http.ListenAndServe(":8080", nil) // nil - null(none)
}

func main() {
	handleRequest()

}
