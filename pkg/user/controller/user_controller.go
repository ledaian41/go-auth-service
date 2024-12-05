package user_controller

import (
	"auth/pkg/user/handler"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.GET("/users", user_handler.GetUserList)
}
