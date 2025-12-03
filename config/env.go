package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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
	TcpPort   string
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
		TcpPort:   os.Getenv("TCP_PORT"),
	}
}
