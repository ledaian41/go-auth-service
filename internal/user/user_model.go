package user

import (
	"go-auth-service/internal/shared"
	"strings"
)

type User struct {
	ID           string `gorm:"primary_key" json:"id"`
	Username     string `gorm:"uniqueIndex:idx_site_username" json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	DOB          string `json:"dob"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `gorm:"uniqueIndex:idx_site_email" json:"email"`
	Avatar       string `json:"avatar"`
	Role         string `json:"role"`
	Site         string `gorm:"uniqueIndex:idx_site_username;uniqueIndex:idx_site_email" json:"site"`
	TokenVersion int    `gorm:"default:0" json:"token_version"`
}

type UserResponse struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	DOB         string `json:"dob"`
	PhoneNumber string `json:"phone"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	Role        string `json:"role"`
}

func (user User) Response() UserResponse {
	return UserResponse{
		Username:    user.Username,
		Name:        user.Name,
		DOB:         user.DOB,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Avatar:      user.Avatar,
		Role:        user.Role,
	}
}

func (user User) ToDTO() shared.UserDTO {
	return shared.UserDTO{
		Username:    user.Username,
		Password:    user.Password,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        user.Role,
		Site:        user.Site,
	}
}

func (user User) HasRole(role string) bool {
	return strings.Contains(user.Role, role)
}

func (user User) HasAnyRole(roles []string) bool {
	for _, requiredRole := range roles {
		if strings.Contains(user.Role, requiredRole) {
			return true
		}
	}
	return false
}
