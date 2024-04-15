package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

const CONFIG_PATH string = "/config/local.yaml"

// Структуры копирующие файл local.yaml
type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-default:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"time_out" env-default:"4s"`
	IdleTimeOut time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	//Получение конфига из переменной окружения
	configPath := os.Getenv(CONFIG_PATH)
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	//проверка существования такого файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)

	}

	return &cfg
}
