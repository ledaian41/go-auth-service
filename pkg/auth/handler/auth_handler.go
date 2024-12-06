package auth_handler

import (
	"auth/pkg/auth/model"
	"auth/pkg/auth/service"
	"auth/pkg/auth/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	token, err := auth_service.GenerateJwtToken(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "jwt create failed"})
		return
	}

	utils.SetCookieToken(c, token)
	c.JSON(http.StatusOK, gin.H{"message": "register success", "token": token})
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

	token, err := auth_service.GenerateJwtToken(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "jwt create failed"})
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
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	mapClaims := claims.(jwt.MapClaims)
	c.JSON(http.StatusOK, gin.H{
		"username": mapClaims["user"].(string),
		"role":     utils.ToStringSlice(mapClaims["role"].([]interface{})),
		"name":     mapClaims["name"].(string),
		"email":    mapClaims["email"].(string),
		"phone":    mapClaims["phone"].(string),
	})
}
