package auth_handler

import (
	"auth/pkg/auth/model"
	"auth/pkg/auth/service"
	"auth/pkg/auth/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	var newAccount auth_model.RegisterAccount
	if err := c.ShouldBind(&newAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := auth_service.CreateNewAccount(&newAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateJwtToken(user.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JWT create failed"})
		return
	}

	utils.SetCookieToken(c, token)
	c.JSON(http.StatusOK, gin.H{"message": "Register Success", "user": user.Response()})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user, err := auth_service.CheckValidUser(username, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateJwtToken(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JWT create failed"})
		return
	}

	utils.SetCookieToken(c, token)
	c.IndentedJSON(http.StatusOK, user.Response())
}

func Logout(c *gin.Context) {
	c.SetCookie("jwt", "", 0, "", "", false, true)
	c.Status(http.StatusOK)
}

func JWT(c *gin.Context) {
	// Get token from Cookie
	tokenStr, err := c.Cookie("jwt")
	if err != nil || len(tokenStr) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No jwt token"})
		return
	}

	claims, err := utils.ExtractJwtToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse JWT token"})
		return
	}

	// Extract user from token
	userId := claims["user"]
	if userId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	user, err := auth_service.GetSessionUser(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user.Response())
}
