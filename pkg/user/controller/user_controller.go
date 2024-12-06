package user_controller

import (
	"auth/middleware"
	"auth/pkg/user/handler"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.GET("/:siteId/users", middleware.AuthMiddleware(), user_handler.GetUserList)
}
