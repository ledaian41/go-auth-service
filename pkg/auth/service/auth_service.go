package auth_service

import (
	"auth/pkg/auth/model"
	"auth/pkg/site/model"
	"auth/pkg/user/model"
	"auth/pkg/user/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

func CheckValidUser(username string, password string) (*user_model.User, error) {
	if len(strings.Trim(username, " ")) == 0 {
		return nil, errors.New("invalid username or password")
	}

	if len(strings.Trim(password, " ")) == 0 {
		return nil, errors.New("invalid username or password")
	}

	user, err := user_model.GetById(username)
	if err != nil {
		return nil, err
	}

	if !CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateNewAccount(account *auth_model.RegisterAccount) (*user_model.User, error) {
	hashedPassword, err := HashPassword(account.Password)
	if err != nil {
		return nil, err
	}

	newUser := user_model.User{
		Username:    account.Username,
		Password:    hashedPassword,
		PhoneNumber: account.PhoneNumber,
		Email:       account.Email,
		Role:        []string{"user"},
	}
	return user_service.CreateNewUser(&newUser)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func GenerateJwtToken(c *gin.Context, user *user_model.User) (string, error) {
	secretKey, err := SecretKeyBySite(c)
	if err != nil {
		return "", err
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":  user.Username,
		"role":  user.Role,
		"name":  user.Name,
		"email": user.Email,
		"phone": user.PhoneNumber,
		"ttl":   time.Now().Add(time.Minute * 30).Unix(), // 30 minutes
	})

	// Sign, get the complete encoded token as a string
	return token.SignedString([]byte(secretKey))
}

func ExtractJwtToken(c *gin.Context, tokenStr string) (jwt.MapClaims, error) {
	secretKey, err := SecretKeyBySite(c)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	// Check expiry of token
	if claims["ttl"].(float64) < float64(time.Now().Unix()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func SecretKeyBySite(c *gin.Context) (string, error) {
	// Check site from middleware
	site, exists := c.Get("site")
	if !exists {
		return "", errors.New("no site info")
	}

	// Check secret key
	secretKey := site.(*site_model.Site).SecretKey
	if len(secretKey) == 0 {
		return "", errors.New("site has no secret key")
	}

	return secretKey, nil
}
