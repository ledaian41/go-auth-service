package main

import (
	"fmt"
	"go-auth-service/config"
	"go-auth-service/internal/auth"
	"go-auth-service/internal/site"
	"go-auth-service/internal/token"
	"go-auth-service/internal/user"
	proto "go-auth-service/proto"
	"go-auth-service/routes"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func startGrpc() {
	db := config.InitDatabase()
	redisClient := config.InitRedisClient()

	tcpPort := os.Getenv("TCP_PORT")
	lis, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		log.Println(fmt.Sprint("❌ Failed to listen:", err))
	}

	grpcServer := grpc.NewServer()

	siteService := site.NewSiteService()
	userService := user.NewUserService(db)
	userService.MigrateDatabase()
	_ = userService.SeedUsers("./internal/user/userData.json")
	tokenService := token.NewTokenService(db)
	tokenService.MigrateDatabase()
	authService := auth.NewAuthService(redisClient, userService, tokenService)
	grpcHandler := auth.New(siteService, authService, tokenService)

	proto.RegisterAuthServer(grpcServer, grpcHandler)
	log.Println("🌤️gRPC server running on :", tcpPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("Failed to serve:", err)
	}
}

func startHttp() {
	db := config.InitDatabase()
	redisClient := config.InitRedisClient()

	r := routes.SetupRouter(db, redisClient)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("🌤️Server running on port", port)
	r.Run(":" + port)
}
