package auth

import (
	"errors"
	"fmt"
	"go-auth-service/config"
	"go-auth-service/internal/shared"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *Service) GenerateRefreshToken(user *shared.UserDTO, sessionId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":          user.Username,
		"site":          user.Site,
		"exp":           jwt.NewNumericDate(time.Now().Add(config.RefreshTokenExpire)),
		"token_version": user.TokenVersion,
		"jti":           sessionId,
	})
	// Sign, get the complete encoded token as a string
	return token.SignedString([]byte(config.Env.SecretKey))
}

func (s *Service) ValidateRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Check JTI
	jti, ok := claims["jti"].(string)
	if !ok {
		return nil, errors.New("invalid jti claim")
	}

	if s.tokenService.ValidateRefreshToken(jti) == "" {
		return nil, errors.New("refresh token revoked or invalid")
	}

	tokenVersion := s.redisClient.GetTokenVersion(claims["user"].(string), claims["site"].(string))
	if int(claims["token_version"].(float64)) < tokenVersion {
		return nil, errors.New("token version expired")
	}

	return claims, nil
}

func (s *Service) GenerateAccessToken(siteSecretKey, sessionId string, user *shared.UserDTO) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":          user.Username,
		"role":          user.Role,
		"name":          user.Name,
		"email":         user.Email,
		"phone":         user.PhoneNumber,
		"exp":           jwt.NewNumericDate(time.Now().Add(config.AccessTokenExpire)),
		"session_id":    sessionId,
		"token_version": user.TokenVersion,
	})

	// Sign, get the complete encoded token as a string
	return token.SignedString([]byte(siteSecretKey))
}

func (s *Service) ValidateAccessToken(site *shared.SiteDTO, tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(site.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	tokenVersion := s.redisClient.GetTokenVersion(claims["user"].(string), site.ID)
	if int(claims["token_version"].(float64)) < tokenVersion {
		return nil, errors.New("token version expired")
	}

	sessionId := claims["session_id"].(string)
	if !s.redisClient.ValidateSessionId(sessionId) {
		return nil, errors.New("token session expired")
	}

	return claims, nil
}

func (s *Service) RevokeSessionId(sessionId string) {
	s.redisClient.AddSessionIdToBlackList(sessionId)
}
