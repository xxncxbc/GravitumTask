package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type DbConfig struct {
	Dsn string
}

type Config struct {
	Db DbConfig
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
	}
}
