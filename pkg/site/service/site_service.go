package site_service

import (
	"go-auth-service/pkg/shared/dto"
	"go-auth-service/pkg/site/data"
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
