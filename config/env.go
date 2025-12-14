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
	RedisPwd  string
	TcpPort   string
	DbHost    string
	DbUser    string
	DbPwd     string
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
		RedisPwd:  os.Getenv("REDIS_PASSWORD"),
		TcpPort:   os.Getenv("TCP_PORT"),
		DbHost:    os.Getenv("DATABASE_HOST"),
		DbUser:    os.Getenv("DATABASE_USER"),
		DbPwd:     os.Getenv("DATABASE_PWD"),
	}
}
