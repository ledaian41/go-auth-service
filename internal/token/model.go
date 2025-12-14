package token

import "time"

type UserToken struct {
	ID        string    `json:"id" gorm:"primary_key"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
