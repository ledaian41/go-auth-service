package auth_service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-auth-service/config"
	"go-auth-service/pkg/shared/dto"
	"time"
)

func (s *AuthService) GenerateAccessToken(siteSecretKey, sessionId string, user *shared_dto.UserDTO) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":          user.Username,
		"role":          user.Role,
		"name":          user.Name,
		"email":         user.Email,
		"phone":         user.PhoneNumber,
		"ttl":           time.Now().Add(accessExpireTime).Unix(),
		"session_id":    sessionId,
		"token_version": user.TokenVersion,
	})

	// Sign, get the complete encoded token as a string
	return token.SignedString([]byte(siteSecretKey))
}

func (s *AuthService) ValidateAccessToken(site *shared_dto.SiteDTO, tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(site.SecretKey), nil
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

	sessionVersion := s.redisClient.GetSessionVersion(claims["user"].(string), site.ID)
	if int(claims["token_version"].(float64)) < sessionVersion {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func (s *AuthService) GenerateRefreshToken(user *shared_dto.UserDTO) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":          user.Username,
		"site":          user.Site,
		"ttl":           time.Now().Add(refreshExpireTime).Unix(),
		"token_version": user.TokenVersion,
	})
	// Sign, get the complete encoded token as a string
	return token.SignedString([]byte(config.Env.SecretKey))
}

func (s *AuthService) ValidateRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.SecretKey), nil
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

	sessionVersion := s.redisClient.GetSessionVersion(claims["user"].(string), claims["site"].(string))
	if int(claims["token_version"].(float64)) < sessionVersion {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
