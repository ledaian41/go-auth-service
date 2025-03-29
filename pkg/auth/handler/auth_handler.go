package auth_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-auth-service/pkg/auth/utils"
	"go-auth-service/pkg/shared/dto"
	"go-auth-service/pkg/shared/interface"
	"net/http"
)

type AuthHandler struct {
	authService shared_interface.AuthServiceInterface
}

func NewAuthHandler(authService shared_interface.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (handler *AuthHandler) Register(c *gin.Context) {
	var newAccount shared_dto.RegisterRequestDTO
	if err := c.ShouldBind(&newAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := handler.authService.CreateNewAccount(&newAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	accessToken, err := handler.authService.GenerateAccessToken(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate access token failed"})
		return
	}

	refreshToken, err := handler.authService.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate refresh token failed"})
		return
	}

	auth_utils.SetCookieToken(c, refreshToken)
	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (handler *AuthHandler) Login(c *gin.Context) {
	var loginAccount shared_dto.LoginRequestDTO
	if err := c.ShouldBind(&loginAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	siteId := c.Param("siteId")
	user, err := handler.authService.CheckValidUser(loginAccount.Username, loginAccount.Password, siteId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("login", user)
	accessToken, err := handler.authService.GenerateAccessToken(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate access token failed"})
		return
	}

	refreshToken, err := handler.authService.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate refresh token failed"})
		return
	}

	auth_utils.SetCookieToken(c, refreshToken)
	c.IndentedJSON(http.StatusOK, gin.H{"token": accessToken})
}

func (handler *AuthHandler) Logout(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	mapClaims := claims.(jwt.MapClaims)
	siteId := c.Param("siteId")
	handler.authService.RevokeUserSession(mapClaims["user"].(string), siteId)
	auth_utils.DestroyCookieToken(c)
	c.String(http.StatusOK, "Signed out")
}

func (handler *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := auth_utils.GetCookieToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "no refresh token"})
		return
	}

	claims, err := handler.authService.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	siteObj, _ := c.Get("site")
	site := siteObj.(*shared_dto.SiteDTO)
	if site.ID != claims["site"].(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "site not matched"})
		return
	}

	user, err := handler.authService.FindUserByUsername(claims["user"].(string), site.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	accessToken, err := handler.authService.GenerateAccessToken(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate access token failed"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": accessToken})
}

func (handler *AuthHandler) JWT(c *gin.Context) {
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
		role = auth_utils.ToStringSlice(role.([]interface{}))
	}

	c.JSON(http.StatusOK, gin.H{
		"username": mapClaims["user"].(string),
		"role":     role,
		"name":     mapClaims["name"].(string),
		"email":    mapClaims["email"].(string),
		"phone":    mapClaims["phone"].(string),
	})
}
