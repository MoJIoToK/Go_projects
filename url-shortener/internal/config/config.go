package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

//const CONFIG_PATH string = "./config/local.yaml"

// region Структуры копирующие настройки из файла local.yaml

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"time_out" env-default:"4s"`
	IdleTimeOut time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

//endregion

// MustLoad - при ошибке должна выдавать панику.
func MustLoad() *Config {
	//Искусственно введенная переменная окружения. Пока не разобрался как сделать по другому
	err := os.Setenv("CONFIG_PATH", "./config/local.yaml")
	if err != nil {
		log.Fatal("Error setting CONFIG_PATH")
	}

	//Получение пути до файла с настройками из переменной окружения
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	//Проверка существования файла с настройками
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return &cfg
}
