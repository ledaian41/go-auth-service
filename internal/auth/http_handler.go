package auth

import (
	"go-auth-service/internal/shared/dto"
	"go-auth-service/internal/shared/interface"
	"go-auth-service/internal/shared/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type HttpHandler struct {
	authService  shared_interface.AuthService
	tokenService shared_interface.TokenService
}

func NewAuthHandler(authService shared_interface.AuthService, tokenService shared_interface.TokenService) *HttpHandler {
	return &HttpHandler{authService: authService, tokenService: tokenService}
}

func (handler *HttpHandler) Register(c *gin.Context) {
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

	refreshToken, err := handler.authService.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate refresh token failed"})
		return
	}
	SetCookieToken(c, refreshToken)
	sessionId := handler.tokenService.StoreRefreshToken(user.Username, refreshToken)
	if len(sessionId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error saving refresh token"})
		return
	}

	site, _ := shared_utils.ReadSiteContext(c)
	accessToken, err := handler.authService.GenerateAccessToken(site.SecretKey, sessionId, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate access token failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (handler *HttpHandler) Login(c *gin.Context) {
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

	refreshToken, err := handler.authService.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "generate refresh token failed"})
		return
	}
	SetCookieToken(c, refreshToken)
	sessionId := handler.tokenService.StoreRefreshToken(user.Username, refreshToken)
	if len(sessionId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error saving refresh token"})
		return
	}

	site, _ := shared_utils.ReadSiteContext(c)
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
	sessionId := handler.tokenService.RevokeRefreshToken(refreshToken)
	if len(sessionId) > 0 {
		handler.authService.RevokeSessionId(sessionId)
	}
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

	sessionId := handler.tokenService.ValidateRefreshToken(refreshToken)
	if len(sessionId) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid refresh token"})
		return
	}

	site, _ := shared_utils.ReadSiteContext(c)
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
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
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
