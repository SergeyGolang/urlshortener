package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-requered:"True"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Addres      string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	Idletimeout time.Duration `yaml:"idle_timeout " env-default:"60s"`
	User        string        `yaml:"user" env-requered:"true"`
	Password    string        `yaml:"password" env-requered:"true" env:"HTTP_SERVER_PASSWORD"`
}

// MustLoad loads configuration from file into Config struct.
// Panics if config cannot be loaded (missing path, invalid file or parsing error).
func MustLoad() *Config {
	// Get config path from environment variable (required)
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// Verify config file exists before attempting to read
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	// Parse config file into struct using cleanenv
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg

}
