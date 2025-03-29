package site_model

import "go-auth-service/pkg/shared/dto"

type Site struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SecretKey string `json:"secret_key"`
}

func (site *Site) ToDTO() shared_dto.SiteDTO {
	return shared_dto.SiteDTO{
		ID:        site.ID,
		SecretKey: site.SecretKey,
	}
}
