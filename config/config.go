package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type AppConfig struct {
	SecretKey string
}

func LoadConfig() *AppConfig {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println(fmt.Sprintf("❌ Error loading .env file: %v", err))
	}

	// Create and return AppConfig instance
	return &AppConfig{
		SecretKey: os.Getenv("SECRET_KEY"),
	}
}
