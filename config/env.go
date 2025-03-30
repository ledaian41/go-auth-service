package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var Env *AppConfig

type AppConfig struct {
	SecretKey string
	RedisHost string
}

func LoadConfig() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println(fmt.Sprintf("❌ Error loading .env file: %v", err))
	}

	Env = &AppConfig{
		SecretKey: os.Getenv("SECRET_KEY"),
		RedisHost: os.Getenv("REDIS_HOST"),
	}
}
