package token

import "time"

type UserToken struct {
	ID           string    `gorm:"primary_key" json:"id"`
	UserID       string    `json:"user_id"`
	RefreshToken string    `gorm:"unique" json:"refresh_token"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
