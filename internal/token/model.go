package token

import "time"

type UserToken struct {
	ID        string    `json:"id" gorm:"primary_key"`
	SiteID    string    `json:"site_id"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
