package site

import (
	"encoding/json"
	"go-auth-service/internal/shared"
	"log"
	"os"

	"gorm.io/gorm"
)

type SiteService struct {
	db *gorm.DB
}

func NewSiteService(db *gorm.DB) *SiteService {
	return &SiteService{db: db}
}

func (s *SiteService) MigrateDatabase() {
	err := s.db.AutoMigrate(&Site{})
	if err != nil {
		log.Printf("❌ Failed at AutoMigrate: %v", err)
	}
	log.Println("🎉 Site - Database migrate successfully")
}

func (s *SiteService) SeedSites(filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File not found is fine, just skip seeding
	}

	// Read file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var sites []Site
	if err := json.Unmarshal(fileData, &sites); err != nil {
		return err
	}

	for _, site := range sites {
		var count int64
		s.db.Model(&Site{}).Where("id = ?", site.ID).Count(&count)
		if count > 0 {
			// Optionally update existing site?
			// For now, let's assume we just ensure it exists.
			continue
		}

		if err := s.db.Create(&site).Error; err != nil {
			log.Printf("Failed to seed site %s: %v", site.Name, err)
		}
	}
	log.Println("🎉 Site seeding completed")
	return nil
}

func (s *SiteService) CheckSiteExists(siteId string) *shared.SiteDTO {
	var site Site
	if err := s.db.Where("id = ?", siteId).First(&site).Error; err != nil {
		return nil
	}
	siteDto := site.ToDTO()
	return &siteDto
}
