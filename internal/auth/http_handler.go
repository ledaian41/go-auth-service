package auth

import (
	"go-auth-service/internal/shared"
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

	sessionId := shared.RandomID()
	refreshToken, err := handler.authService.GenerateRefreshToken(user, sessionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate refresh token failed"})
		return
	}
	SetCookieToken(c, refreshToken)
	_, err = handler.tokenService.StoreRefreshToken(user.Username, sessionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error saving refresh token"})
		return
	}

	site, _ := shared.ReadSiteContext(c)
	accessToken, err := handler.authService.GenerateAccessToken(site.SecretKey, sessionId, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate access token failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (handler *HttpHandler) Login(c *gin.Context) {
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

	sessionId := shared.RandomID()
	refreshToken, err := handler.authService.GenerateRefreshToken(user, sessionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate refresh token failed"})
		return
	}
	SetCookieToken(c, refreshToken)
	_, err = handler.tokenService.StoreRefreshToken(user.Username, sessionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error saving refresh token"})
		return
	}

	site, _ := shared.ReadSiteContext(c)
	accessToken, err := handler.authService.GenerateAccessToken(site.SecretKey, sessionId, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate access token failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (handler *HttpHandler) Logout(c *gin.Context) {
	refreshToken, err := GetCookieToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "no refresh token"})
		return
	}

	claims, err := handler.authService.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		return
	}

	sessionId := claims["jti"].(string)
	handler.tokenService.RevokeRefreshToken(sessionId)
	handler.authService.RevokeSessionId(sessionId)
	DestroyCookieToken(c)
	c.String(http.StatusOK, "Signed out")
}

func (handler *HttpHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := GetCookieToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "no refresh token"})
		return
	}

	claims, err := handler.authService.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	sessionId := claims["jti"].(string)

	site, _ := shared.ReadSiteContext(c)
	if site.ID != claims["site"].(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "site not matched"})
		return
	}

	user, err := handler.authService.FindUserByUsername(claims["user"].(string), site.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	accessToken, err := handler.authService.GenerateAccessToken(site.SecretKey, sessionId, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate access token failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": accessToken})
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
