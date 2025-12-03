package main

import (
	"fmt"
	"go-auth-service/config"
	user_model "go-auth-service/pkg/user/model"
	auth "go-auth-service/proto"
	"go-auth-service/routes"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func startGrpc() {
	tcpPort := os.Getenv("TCP_PORT")
	lis, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		log.Println(fmt.Sprint("❌ Failed to listen:", err))
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServer(grpcServer, &AuthServer{})
	log.Println("🌤️gRPC server running on :", tcpPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("Failed to serve:", err)
	}
}

func startHttp() {
	db := config.InitDatabase()
	redisClient := config.InitRedisClient()

	_ = user_model.LoadUsersFromFile("./pkg/user/data/userData.json")

	r := routes.SetupRouter(db, redisClient)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("🌤️Server running on port", port)
	r.Run(":" + port)
}
