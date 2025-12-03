package site_service

import (
	"go-auth-service/internal/shared/dto"
	"go-auth-service/internal/site/data"
)

type SiteService struct{}

func NewSiteService() *SiteService {
	return &SiteService{}
}

func (s *SiteService) CheckSiteExists(siteId string) *shared_dto.SiteDTO {
	for _, site := range site_data.SiteData() {
		if site.ID == siteId {
			siteCopy := site.ToDTO()
			return &siteCopy
		}
	}
	return nil
}
