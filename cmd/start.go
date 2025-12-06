package main

import (
	"context"
	"fmt"
	"go-auth-service/config"
	"go-auth-service/internal/user"
	auth "go-auth-service/proto"
	"go-auth-service/routes"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func (s *AuthServer) Ping(ctx context.Context, req *auth.HelloRequest) (*auth.HelloResponse, error) {
	log.Printf("Received: %s", req.Name)
	return &auth.HelloResponse{Message: "Hello, " + req.Name}, nil
}

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

	_ = user.LoadUsersFromFile("./internal/user/userData.json")

	r := routes.SetupRouter(db, redisClient)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("🌤️Server running on port", port)
	r.Run(":" + port)
}
