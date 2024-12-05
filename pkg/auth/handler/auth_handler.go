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
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := auth_service.CreateNewAccount(&newAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenerateJwtToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "JWT create failed"})
		return
	}

	utils.SetCookieToken(c, token)
	c.JSON(http.StatusOK, gin.H{"message": "Register Success", "token": token})
}

func Login(c *gin.Context) {
	var loginAccount auth_model.LoginAccount
	if err := c.ShouldBind(&loginAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := auth_service.CheckValidUser(loginAccount.Username, loginAccount.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenerateJwtToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "JWT create failed"})
		return
	}

	utils.SetCookieToken(c, token)
	c.IndentedJSON(http.StatusOK, gin.H{"token": token})
}

func Logout(c *gin.Context) {
	utils.DestroyCookieToken(c)
	c.Status(http.StatusOK)
}

func JWT(c *gin.Context) {
	// Get token from Cookie
	tokenStr, err := utils.GetCookieToken(c)
	if err != nil || len(tokenStr) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No jwt token"})
		return
	}

	claims, err := utils.ExtractJwtToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed to parse JWT token"})
		return
	}

	// Extract user from token
	userId := claims["user"]
	if userId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": claims["user"].(string),
		"email":    claims["email"].(string),
		"role":     claims["role"].(string),
		"site":     claims["site"].(string),
	})
}
