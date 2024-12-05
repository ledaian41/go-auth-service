package site_service

import (
	"auth/pkg/site/data"
	"auth/pkg/site/model"
)

func CheckSiteExists(siteId string) *site_model.Site {
	for _, site := range site_data.SiteData() {
		if site.ID == siteId {
			siteCopy := site
			return &siteCopy
		}
	}
	return nil
}
