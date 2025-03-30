package token_service

import (
	"go-auth-service/pkg/shared/utils"
	"go-auth-service/pkg/token/model"
	"gorm.io/gorm"
	"log"
)

type TokenService struct {
	db *gorm.DB
}

func NewTokenService(db *gorm.DB) *TokenService {
	return &TokenService{db: db}
}

func (s *TokenService) MigrateDatabase() {
	err := s.db.AutoMigrate(&token_model.UserToken{})
	if err != nil {
		log.Printf("❌ Failed at AutoMigrate: %v", err)
	}
	log.Println("🎉 NodeType - Database migrate successfully")
}

func (s *TokenService) ValidateRefreshToken(refreshToken string) string {
	var token token_model.UserToken
	_ = s.db.Where("refresh_token = ?", refreshToken).First(&token).Error
	return token.ID
}

func (s *TokenService) StoreRefreshToken(username, refreshToken string) string {
	sessionId := shared_utils.RandomID()
	token := token_model.UserToken{
		ID:           sessionId,
		UserID:       username,
		RefreshToken: refreshToken,
	}
	if err := s.db.Create(&token).Error; err != nil {
		log.Println("error create token in database", err)
		return ""
	}
	return sessionId
}

func (s *TokenService) RevokeRefreshToken(refreshToken string) string {
	var token token_model.UserToken
	if err := s.db.Where("refresh_token = ?", refreshToken).First(&token).Error; err != nil {
		return ""
	}

	s.db.Where("refresh_token = ?", refreshToken).Delete(&token_model.UserToken{})
	return token.ID
}
