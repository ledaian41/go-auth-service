package shared

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Filter[T any](arr []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range arr {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, R any](input []T, mapper func(T) R) []R {
	var result []R
	for _, item := range input {
		result = append(result, mapper(item))
	}
	return result
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RandomID() string {
	b := make([]byte, 4) // 4 bytes = 8 hex characters
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("❌ Failed to generate random ID: %v", err)
	}
	return hex.EncodeToString(b)
}

func ReadSiteContext(c *gin.Context) (*SiteDTO, error) {
	// Check site from middleware
	site, exists := c.Get("site")
	if !exists {
		return nil, errors.New("no site info, have to use site middleware")
	}

	return site.(*SiteDTO), nil
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
