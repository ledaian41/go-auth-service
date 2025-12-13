package site

import (
	"go-auth-service/internal/shared"
)

type SiteService struct{}

func NewSiteService() *SiteService {
	return &SiteService{}
}

func (s *SiteService) CheckSiteExists(siteId string) *shared.SiteDTO {
	for _, site := range GetData() {
		if site.ID == siteId {
			siteCopy := site.ToDTO()
			return &siteCopy
		}
	}
	return nil
}
