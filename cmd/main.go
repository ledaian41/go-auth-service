package main

import (
	_ "auth/docs"
	"auth/middleware"
	auth_controller "auth/pkg/auth/controller"
	user_controller "auth/pkg/user/controller"
	user_model "auth/pkg/user/model"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
)

// @title Authentication Service API
// @version 1.0
// @description The Core Authentication Service is a microservice designed to handle user authentication and provide JWT tokens for secure access. Third-party applications can integrate with this service to authenticate users and validate their identities.
// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()

	_ = user_model.LoadUsersFromFile("./pkg/user/data/userData.json")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(middleware.SiteMiddleware())

	go auth_controller.Router(r)

	go user_controller.Router(r)

	r.Run("localhost:8080")
}
