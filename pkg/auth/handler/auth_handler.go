package auth_handler

import (
	auth_model "auth/pkg/auth/model"
	"auth/pkg/auth/service"
	"auth/pkg/auth/utils"
	site_model "auth/pkg/site/model"
	user_model "auth/pkg/user/model"
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

	accessToken, err := auth_service.GenerateAccessToken(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate access token failed"})
		return
	}

	refreshToken, err := auth_service.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate refresh token failed"})
		return
	}

	utils.SetCookieToken(c, refreshToken)
	c.JSON(http.StatusOK, gin.H{"token": accessToken})
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

	accessToken, err := auth_service.GenerateAccessToken(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate access token failed"})
		return
	}

	refreshToken, err := auth_service.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate refresh token failed"})
		return
	}

	utils.SetCookieToken(c, refreshToken)
	c.IndentedJSON(http.StatusOK, gin.H{"token": accessToken})
}

func Logout(c *gin.Context) {
	utils.DestroyCookieToken(c)
	c.Status(http.StatusOK)
}

func RefreshToken(c *gin.Context) {
	refreshToken, err := utils.GetCookieToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "no refresh token"})
		return
	}

	claims, err := auth_service.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	siteObj, _ := c.Get("site")
	site := siteObj.(*site_model.Site)
	if site.ID != claims["site"].(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "site not matched"})
		return
	}

	user, err := user_model.GetById(claims["user"].(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	accessToken, err := auth_service.GenerateAccessToken(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate access token failed"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": accessToken})
}

func JWT(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	mapClaims := claims.(jwt.MapClaims)
	role := mapClaims["role"]
	if role == nil {
		role = ""
	} else {
		role = utils.ToStringSlice(role.([]interface{}))
	}

	c.JSON(http.StatusOK, gin.H{
		"username": mapClaims["user"].(string),
		"role":     role,
		"name":     mapClaims["name"].(string),
		"email":    mapClaims["email"].(string),
		"phone":    mapClaims["phone"].(string),
	})
}
