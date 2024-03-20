package main

import (
	"log"
	"middleware/internal/pkg/app"
)

func main() {

	//Creating an Instance of App.
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	//Starting the App.
	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}

}
