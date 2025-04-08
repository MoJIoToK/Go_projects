package main

import (
	"GoNews/internal/api"
	"GoNews/internal/config"
	"GoNews/internal/logger"
	"GoNews/internal/model"
	"GoNews/internal/rss"
	"GoNews/internal/storage"
	"GoNews/internal/storage/mongo"
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"
)

// server - структура сервера, хранящая экземпляр БД и api.
type server struct {
	db  storage.DB
	api *api.API
}

func main() {

	//Инициализация логера.
	logger.SetupLogger()
	slog.Debug("Logger setup load successful")

	var srv server

	//Загрузка конфигураций из файла конфигурации.
	cfg := config.MustLoad("./config/config.yaml")
	slog.Debug("Load config file success")

	//Инициализация зависимостей
	conn := cfg.StoragePath
	db, err := mongo.New(conn)
	if err != nil {
		log.Fatal(err)
	}

	srv.db = db
	srv.api = api.New(srv.db)

	//Создание каналов для публикаций и ошибок
	chPosts := make(chan []model.Post)
	chErrors := make(chan error)

	//Запуск парсинга новостей в отдельном потоке для каждой ссылки
	for _, url := range cfg.URLS {
		go parseURL(url, chPosts, chErrors, cfg.Period)
	}

	ctx := context.Background()
	//Запись потока публикаций из канала в БД
	go func() {
		for posts := range chPosts {
			db.AddPost(ctx, posts)
		}
	}()

	//Обработка потока ошибок из канала
	go func() {
		for err := range chErrors {
			log.Println("Error:", err)
		}
	}()

	//Запуск веб-сервера с API и приложением
	err = http.ListenAndServe(cfg.Address, srv.api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

// parseURL - функция позволяет асинхронно читать RSS-поток. Раскодированные новости и ошибки
// записываются в соответствующие каналы
func parseURL(url string, posts chan []model.Post, errs chan error, period int) {
	for {
		news, err := rss.Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
	}
	time.Sleep(time.Minute * time.Duration(period))
}
