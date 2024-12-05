package main

import (
	"auth/middleware"
	"auth/pkg/auth/controller"
	"auth/pkg/user/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.Use(middleware.SiteMiddleware())

	go auth_controller.Router(r)

	go user_controller.Router(r)

	r.Run("localhost:8080")
}
