package auth_service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-auth-service/config"
	"go-auth-service/pkg/shared/dto"
	"go-auth-service/pkg/shared/interface"
	"go-auth-service/pkg/shared/utils"
	"strings"
	"time"
)

var (
	accessExpireTime  = time.Minute * 15   // 15 minutes
	refreshExpireTime = time.Hour * 24 * 7 // 1 week
)

type AuthService struct {
	userService shared_interface.UserServiceInterface
	secretKey   string
}

func NewAuthService(userService shared_interface.UserServiceInterface, secretKey string) *AuthService {
	return &AuthService{
		userService: userService,
		secretKey:   secretKey,
	}
}

func (s *AuthService) CheckValidUser(username, password, siteId string) (*shared_dto.UserDTO, error) {
	if len(strings.Trim(username, " ")) == 0 {
		return nil, errors.New("invalid username or password")
	}

	if len(strings.Trim(password, " ")) == 0 {
		return nil, errors.New("invalid username or password")
	}

	user, err := s.FindUserByUsername(username, siteId)
	if err != nil {
		return nil, err
	}

	if !shared_utils.CheckHashPassword(password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *AuthService) CreateNewAccount(account *shared_dto.RegisterRequestDTO) (*shared_dto.UserDTO, error) {
	hashedPassword, err := shared_utils.HashPassword(account.Password)
	if err != nil {
		return nil, err
	}

	newUser := shared_dto.UserDTO{
		Username:    account.Username,
		Password:    hashedPassword,
		PhoneNumber: account.PhoneNumber,
		Email:       account.Email,
		Role:        []string{"user"},
	}
	return s.userService.CreateNewUser(&newUser)
}

func (s *AuthService) GenerateAccessToken(c *gin.Context, user *shared_dto.UserDTO) (string, error) {
	siteSecretKey, err := secretKeyBySite(c)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":          user.Username,
		"role":          user.Role,
		"name":          user.Name,
		"email":         user.Email,
		"phone":         user.PhoneNumber,
		"ttl":           time.Now().Add(accessExpireTime).Unix(),
		"token_version": user.TokenVersion,
	})

	// Sign, get the complete encoded token as a string
	return token.SignedString([]byte(siteSecretKey))
}

func (s *AuthService) ValidateAccessToken(c *gin.Context, tokenStr string) (jwt.MapClaims, error) {
	siteSecretKey, err := secretKeyBySite(c)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(siteSecretKey), nil
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

func (s *AuthService) GenerateRefreshToken(user *shared_dto.UserDTO) (string, error) {
	appConfig := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.Username,
		"site": user.Site,
		"ttl":  time.Now().Add(refreshExpireTime).Unix(),
	})

	// Sign, get the complete encoded token as a string
	return token.SignedString([]byte(appConfig.SecretKey))
}

func (s *AuthService) ValidateRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	appConfig := config.LoadConfig()
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(appConfig.SecretKey), nil
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

func (s *AuthService) FindUserByUsername(username, siteId string) (*shared_dto.UserDTO, error) {
	return s.userService.FindUserByUsername(username, siteId)
}

func (s *AuthService) CheckAdminRole(role []interface{}) bool {
	for _, r := range role {
		if r == "admin" {
			return true
		}
	}
	return false
}

func secretKeyBySite(c *gin.Context) (string, error) {
	// Check site from middleware
	site, exists := c.Get("site")
	if !exists {
		return "", errors.New("no site info")
	}

	// Check secret key
	secretKey := site.(*shared_dto.SiteDTO).SecretKey
	if len(secretKey) == 0 {
		return "", errors.New("site has no secret key")
	}

	return secretKey, nil
}

func (s *AuthService) RevokeUserSession(username string) {

}
