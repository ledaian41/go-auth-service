package main

import (
	"go-auth-service/config"
	"go-auth-service/internal/site"
	"go-auth-service/internal/token"
	"go-auth-service/internal/user"
)

// @title Authentication Service API
// @version 1.0
// @description The Core Authentication Service is a microservice designed to handle user authentication and provide JWT tokens for secure access. Third-party applications can integrate with this service to authenticate users and validate their identities.
// @host localhost:8080
// @BasePath /
func main() {
	config.LoadConfig()

	db := config.InitDatabase(config.Env.DbHost, config.Env.DbUser, config.Env.DbPwd)
	redisClient := config.InitRedisClient()

	// Migrations and Seeding
	siteService := site.NewSiteService(db)
	siteService.MigrateDatabase()
	_ = siteService.SeedSites("./internal/site/siteData.json")

	userService := user.NewUserService(db)
	userService.MigrateDatabase()
	_ = userService.SeedUsers("./internal/user/userData.json")

	tokenService := token.NewTokenService(db)
	tokenService.MigrateDatabase()

	go startHttp(db, redisClient) // Run Gin HTTP server

	go startGrpc(db, redisClient) // Run gRPC server

	select {}
}
