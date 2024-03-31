package main

import (
	tgClient "article-advisor/clients/telegram"
	"article-advisor/consumer/event-consumer"
	"article-advisor/events/telegram"
	"article-advisor/storage/files"
	"flag"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
	//sqliteStoragePath = "files_storage"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {

	//token := mustToken()

	//tgClient := telegram.New(tgBotHost, mustToken())

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	//fetcher = fetcher.New()
	//
	//processor = processor.New()

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
