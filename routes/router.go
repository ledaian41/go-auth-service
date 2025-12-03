package routes

import (
	"go-auth-service/config"
	_ "go-auth-service/docs"
	"go-auth-service/internal/auth/handler"
	"go-auth-service/internal/auth/service"
	"go-auth-service/internal/site/service"
	"go-auth-service/internal/token/service"
	"go-auth-service/internal/user/handler"
	"go-auth-service/internal/user/service"
	"go-auth-service/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, redisClient *config.RedisClient) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	siteService := site_service.NewSiteService()
	r.Use(middleware.SiteMiddleware(siteService))

	userService := user_service.NewUserService(db)
	userService.MigrateDatabase()

	tokenService := token_service.NewTokenService(db)
	tokenService.MigrateDatabase()

	authService := auth_service.NewAuthService(redisClient, userService)

	AuthHandler := auth_handler.NewAuthHandler(authService, tokenService)
	r.GET("/:siteId/jwt", middleware.AuthMiddleware(authService), AuthHandler.JWT)
	r.GET("/:siteId/refresh", AuthHandler.RefreshToken)
	r.POST("/:siteId/signup", AuthHandler.Register)
	r.POST("/:siteId/login", AuthHandler.Login)
	r.GET("/:siteId/signout", AuthHandler.Logout)

	UserHandler := user_handler.NewUserHandler(userService)
	r.GET("/:siteId/users", middleware.AuthMiddleware(authService), middleware.AdminAuthMiddleware(authService), UserHandler.GetUserList)

	return r
}
