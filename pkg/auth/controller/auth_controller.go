package auth_controller

import (
	"auth/pkg/auth/handler"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.GET("/jwt", auth_handler.JWT)
	r.POST("/signup", auth_handler.Register)
	r.POST("/login", auth_handler.Login)
	r.GET("/signout", auth_handler.Logout)
}
