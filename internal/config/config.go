package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-requered:"True"` //env-requered-чтобы приложение не запустилось если нет srotage path
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Addres      string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	Idletimeout time.Duration `yaml:"idle_timeout " env-default:"60s"`
	User        string        `yaml:"user" env-requered:"true"`
	Password    string        `yaml:"password" env-requered:"true" env:"HTTP_SERVER_PASSWORD"`
}

/* Приставка Must значит, что функция вместо возврата ошибки будет паниковать. Инициализация конфига - случай паники*/

// Функция загрузки конфига в структуру
func MustLoad() *Config {
	// загрузка пути конфига из переменной окружения и проверка не равна ли перменная != ""
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// проверка существования файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg

}
