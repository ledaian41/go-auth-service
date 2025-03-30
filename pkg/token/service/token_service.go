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

func (s *TokenService) ValidateRefreshToken(refreshToken string) (bool, error) {
	var token token_model.UserToken

	err := s.db.Where("refresh_token = ? AND revoked = ?", refreshToken, false).
		First(&token).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *TokenService) StoreRefreshToken(username, refreshToken string) error {
	token := token_model.UserToken{
		ID:           shared_utils.RandomID(),
		UserID:       username,
		RefreshToken: refreshToken,
	}
	return s.db.Create(&token).Error
}

func (s *TokenService) RevokeRefreshToken(refreshToken string) error {
	return s.db.Model(&token_model.UserToken{}).
		Where("refresh_token = ?", refreshToken).
		Delete(&token_model.UserToken{}).Error
}
