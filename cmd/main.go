package main

import (
	"auth/internal/auth/controller"
	"auth/internal/user/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	auth_controller.Router(r)

	user_controller.Router(r)

	r.Run("localhost:8080")
}
