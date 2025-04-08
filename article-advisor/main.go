package main

import (
	tgClient "article-advisor/clients/telegram"
	"article-advisor/consumer/event-consumer"
	"article-advisor/events/telegram"
	"article-advisor/storage/sqlite"
	"context"
	"flag"
	"log"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/database.db"
	batchSize         = 100
)

func main() {

	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to database: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init database: ", err)
	}

	//Create telegram.Processor
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	//Start Consumer
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

// mustToken allows you to get a token from a flag. Errors are not returned,
// instead the program crashes (therefore a must).
func mustToken() string {

	//Defines a string flag with specified name. The return value is the address of a string
	//variable that stores the value of the flag.
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telebot",
	)

	//Parses the command-line flags.
	flag.Parse()

	//Program crashes without token or if token = ""(empty string).
	if *token == "" {
		log.Fatal("Token is not specified!")
	}
	return *token
}
