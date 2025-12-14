package token

import (
	"log"

	"gorm.io/gorm"
)

type TokenService struct {
	db *gorm.DB
}

func NewTokenService(db *gorm.DB) *TokenService {
	return &TokenService{db: db}
}

func (s *TokenService) MigrateDatabase() {
	err := s.db.AutoMigrate(&UserToken{})
	if err != nil {
		log.Printf("❌ Failed at AutoMigrate: %v", err)
	}
	log.Println("🎉 NodeType - Database migrate successfully")
}

func (s *TokenService) ValidateRefreshToken(id string) string {
	var token UserToken
	_ = s.db.Where("id = ?", id).First(&token).Error
	return token.ID
}

func (s *TokenService) StoreRefreshToken(id, username, siteId string) (string, error) {
	token := UserToken{
		ID:       id,
		UserName: username,
		SiteID:   siteId,
	}
	if err := s.db.Create(&token).Error; err != nil {
		log.Println("error create token in database", err)
		return "", err
	}
	return id, nil
}

func (s *TokenService) RevokeRefreshToken(id string) string {
	var token UserToken
	if err := s.db.Where("id = ?", id).First(&token).Error; err != nil {
		return ""
	}

	s.db.Where("id = ?", id).Delete(&UserToken{})
	return token.ID
}
