package utils

import (
	"auth/config"
	site_model "auth/pkg/site/model"
	user_model "auth/pkg/user/model"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

func SetCookieToken(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	site, _ := c.Get("site")
	siteId := site.(*site_model.Site).ID
	c.SetCookie("jwt-"+siteId, token, 60*30, "", "", false, true)
}

func GetCookieToken(c *gin.Context) (string, error) {
	site, _ := c.Get("site")
	siteId := site.(*site_model.Site).ID
	return c.Cookie("jwt-" + siteId)
}

func DestroyCookieToken(c *gin.Context) {
	site, _ := c.Get("site")
	siteId := site.(*site_model.Site).ID
	c.SetCookie("jwt-"+siteId, "", 0, "", "", false, true)
}

func GenerateJwtToken(user *user_model.User) (string, error) {
	appConfig := config.LoadConfig()
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":  user.Username,
		"email": user.Email,
		"role":  user.Role,
		"site":  user.Site,
		"ttl":   time.Now().Add(time.Minute * 30).Unix(), // 30 minutes
	})

	// Sign, get the complete encoded token as a string
	return token.SignedString([]byte(appConfig.SecretKey))
}

func ExtractJwtToken(tokenStr string) (jwt.MapClaims, error) {
	appConfig := config.LoadConfig()
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(appConfig.SecretKey), nil
	})
	if err != nil {
		return nil, errors.New("failed to parse JWT token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("JWT Claims failed")
	}

	// Check expiry of token
	if claims["ttl"].(float64) < float64(time.Now().Unix()) {
		return nil, errors.New("Token expired")
	}

	return claims, nil
}
