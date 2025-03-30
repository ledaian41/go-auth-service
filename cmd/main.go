package main

import (
	"go-auth-service/config"
	"go-auth-service/pkg/user/model"
	"go-auth-service/routes"
	"log"
	"os"
)

// @title Authentication Service API
// @version 1.0
// @description The Core Authentication Service is a microservice designed to handle user authentication and provide JWT tokens for secure access. Third-party applications can integrate with this service to authenticate users and validate their identities.
// @host localhost:8080
// @BasePath /
func main() {
	config.LoadConfig()
	redisClient := config.InitRedisClient()

	_ = user_model.LoadUsersFromFile("./pkg/user/data/userData.json")

	r := routes.SetupRouter(redisClient)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	r.Run(":" + port)
}
