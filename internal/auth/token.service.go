package auth

import (
	"errors"
	"fmt"
	"go-auth-service/config"
	"go-auth-service/internal/shared"
	"log"
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

func (s *Service) ParseRefreshToken(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
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

func (s *Service) RotateRefreshToken(site *shared.SiteDTO, oldRefreshToken string) (string, string, error) {
	// 1. Validate old refresh token (format & expiration)
	claims, err := s.ValidateRefreshToken(oldRefreshToken)
	if err != nil {
		return "", "", err
	}

	// 2. Check JTI in Allowlist
	jti, ok := claims["jti"].(string)
	if !ok {
		return "", "", errors.New("invalid jti claim")
	}

	if s.tokenService.ValidateRefreshToken(jti) == "" {
		return "", "", errors.New("refresh token revoked or invalid")
	}

	// 3. Revoke old JTI (One-time usage policy)
	s.tokenService.RevokeRefreshToken(jti)

	// 4. Verify Site Match
	tokenSiteId, ok := claims["site"].(string)
	if !ok || tokenSiteId != site.ID {
		return "", "", errors.New("site mismatch")
	}

	// 5. Find User
	username, ok := claims["user"].(string)
	if !ok {
		return "", "", errors.New("invalid user claim")
	}

	user, err := s.userService.FindUserByUsername(username, site.ID)
	if err != nil {
		return "", "", err
	}

	// 6. Generate New session (JTI)
	newJti := shared.RandomID()

	// 7. Generate New Refresh Token
	newRefreshToken, err := s.GenerateRefreshToken(user, newJti)
	if err != nil {
		return "", "", err
	}

	// 8. Store New JTI
	go func() {
		_, err := s.tokenService.StoreRefreshToken(jti, user.Username, site.ID)
		if err != nil {
			log.Printf("store_refresh_token_failed, jti %s, error: %v", jti, err)
		}
	}()

	// 9. Generate New Access Token
	newAccessToken, err := s.GenerateAccessToken(site.SecretKey, newJti, user)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
