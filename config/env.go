package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var Env *AppConfig

const (
	AccessTokenExpire  = time.Minute * 15   // 15 minutes
	RefreshTokenExpire = time.Hour * 24 * 7 // 1 week
)

type AppConfig struct {
	SecretKey string
	RedisHost string
	CachePath string
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
		CachePath: os.Getenv("CACHE_PATH"),
	}
}
