package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	DSN string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default config")
		return &Config{}
	}

	return &Config{
		DB: DBConfig{
			DSN: os.Getenv("DSN"),
		},
	}
}
