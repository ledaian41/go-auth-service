package auth

import (
	"go-auth-service/internal/shared"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type HttpHandler struct {
	authService  shared.AuthService
	tokenService shared.TokenService
}

func NewAuthHandler(authService shared.AuthService, tokenService shared.TokenService) *HttpHandler {
	return &HttpHandler{authService: authService, tokenService: tokenService}
}

func (handler *HttpHandler) Register(c *gin.Context) {
	site, _ := shared.ReadSiteContext(c)
	var newAccount shared.RegisterRequestDTO
	if err := c.ShouldBind(&newAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := handler.authService.CreateNewAccount(&newAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	jti := shared.RandomID()
	refreshToken, err := handler.authService.GenerateRefreshToken(user, jti)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate refresh token failed"})
		return
	}
	SetCookieToken(c, refreshToken)
	_, err = handler.tokenService.StoreRefreshToken(jti, user.Username, site.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error saving refresh token"})
		return
	}

	accessToken, err := handler.authService.GenerateAccessToken(site.SecretKey, jti, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate access token failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (handler *HttpHandler) Login(c *gin.Context) {
	site, _ := shared.ReadSiteContext(c)
	var loginAccount shared.LoginRequestDTO
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

	jti := shared.RandomID()
	refreshToken, err := handler.authService.GenerateRefreshToken(user, jti)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate refresh token failed"})
		return
	}
	SetCookieToken(c, refreshToken)
	go func() {
		_, err := handler.tokenService.StoreRefreshToken(jti, user.Username, site.ID)
		if err != nil {
			log.Printf("store_refresh_token_failed, jti %s, error: %v", jti, err)
		}
	}()

	accessToken, err := handler.authService.GenerateAccessToken(site.SecretKey, jti, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate access token failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (handler *HttpHandler) Logout(c *gin.Context) {
	refreshToken, err := GetCookieToken(c)
	if err != nil {
		c.String(http.StatusOK, "Signed out")
		return
	}

	claims, err := handler.authService.ParseRefreshToken(refreshToken)
	if err != nil {
		c.String(http.StatusOK, "Signed out")
		return
	}

	sessionId := claims["jti"].(string)
	handler.authService.RevokeSessionId(sessionId)
	go handler.tokenService.RevokeRefreshToken(sessionId)
	DestroyCookieToken(c)
	c.String(http.StatusOK, "Signed out")
}

func (handler *HttpHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := GetCookieToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "no refresh token"})
		return
	}

	site, _ := shared.ReadSiteContext(c)
	if site == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "site context missing"})
		return
	}

	newAccessToken, newRefreshToken, err := handler.authService.RotateRefreshToken(site, refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	SetCookieToken(c, newRefreshToken)
	c.JSON(http.StatusOK, gin.H{"token": newAccessToken})
}

func (handler *HttpHandler) JWT(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user not found"})
		return
	}

	mapClaims := claims.(jwt.MapClaims)

	c.JSON(http.StatusOK, gin.H{
		"username": mapClaims["user"].(string),
		"role":     mapClaims["role"].(string),
		"name":     mapClaims["name"].(string),
		"email":    mapClaims["email"].(string),
		"phone":    mapClaims["phone"].(string),
	})
}
