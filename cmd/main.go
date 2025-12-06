package main

import (
	"go-auth-service/config"
)

// @title Authentication Service API
// @version 1.0
// @description The Core Authentication Service is a microservice designed to handle user authentication and provide JWT tokens for secure access. Third-party applications can integrate with this service to authenticate users and validate their identities.
// @host localhost:8080
// @BasePath /
func main() {
	config.LoadConfig()

	go startHttp() // Run Gin HTTP server

	go startGrpc() // Run gRPC server

	select {}
}
