package site

import "go-auth-service/internal/shared"

type Site struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SecretKey string `json:"secret_key"`
}

func (site *Site) ToDTO() shared.SiteDTO {
	return shared.SiteDTO{
		ID:        site.ID,
		SecretKey: site.SecretKey,
	}
}
