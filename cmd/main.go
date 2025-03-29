package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go-auth-service/config"
	_ "go-auth-service/docs"
	"go-auth-service/middleware"
	"go-auth-service/pkg/auth/handler"
	"go-auth-service/pkg/auth/service"
	"go-auth-service/pkg/site/service"
	"go-auth-service/pkg/user/handler"
	"go-auth-service/pkg/user/model"
	"go-auth-service/pkg/user/service"
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

	siteService := site_service.NewSiteService()
	r.Use(middleware.SiteMiddleware(siteService))

	userService := user_service.NewUserService()
	authService := auth_service.NewAuthService(userService, config.LoadConfig().SecretKey)

	AuthHandler := auth_handler.NewAuthHandler(authService)
	r.GET("/:siteId/jwt", middleware.AuthMiddleware(authService), AuthHandler.JWT)
	r.GET("/:siteId/refresh", AuthHandler.RefreshToken)
	r.POST("/:siteId/signup", AuthHandler.Register)
	r.POST("/:siteId/login", AuthHandler.Login)
	r.GET("/:siteId/signout", middleware.AuthMiddleware(authService), AuthHandler.Logout)

	UserHandler := user_handler.NewUserHandler(userService)
	r.GET("/:siteId/users", middleware.AuthMiddleware(authService), middleware.AdminAuthMiddleware(authService), UserHandler.GetUserList)

	r.Run("localhost:8080")
}
