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
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {

	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	//Запуск Consumer
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

// Получение токена. Ошибки не возвращаются, вместо этого программа аварийно завершается(поэтому must).
func mustToken() string {

	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("Token is not specified")
	}
	return *token
}
