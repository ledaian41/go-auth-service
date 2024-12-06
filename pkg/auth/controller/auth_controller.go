package auth_controller

import (
	"auth/middleware"
	"auth/pkg/auth/handler"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.GET("/:siteId/jwt", middleware.AuthMiddleware(), auth_handler.JWT)
	r.POST("/:siteId/signup", auth_handler.Register)
	r.POST("/:siteId/login", auth_handler.Login)
	r.GET("/:siteId/signout", auth_handler.Logout)
}
