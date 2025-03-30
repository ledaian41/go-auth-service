package routes

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
	"go-auth-service/pkg/user/service"
	"net/http"
)

func SetupRouter(redisClient *config.RedisClient) *gin.Engine {
	r := gin.Default()
	//gin.SetMode(gin.ReleaseMode) for production

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	siteService := site_service.NewSiteService()
	r.Use(middleware.SiteMiddleware(siteService))

	userService := user_service.NewUserService()
	authService := auth_service.NewAuthService(userService, redisClient)

	AuthHandler := auth_handler.NewAuthHandler(authService)
	r.GET("/:siteId/jwt", middleware.AuthMiddleware(authService), AuthHandler.JWT)
	r.GET("/:siteId/refresh", AuthHandler.RefreshToken)
	r.POST("/:siteId/signup", AuthHandler.Register)
	r.POST("/:siteId/login", AuthHandler.Login)
	r.GET("/:siteId/signout", middleware.AuthMiddleware(authService), AuthHandler.Logout)

	UserHandler := user_handler.NewUserHandler(userService)
	r.GET("/:siteId/users", middleware.AuthMiddleware(authService), middleware.AdminAuthMiddleware(authService), UserHandler.GetUserList)

	return r
}
