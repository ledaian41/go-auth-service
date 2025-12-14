package user

import (
	"go-auth-service/internal/shared"
	"strings"
)

type User struct {
	ID           string `json:"id" gorm:"primary_key"`
	Username     string `json:"username" gorm:"uniqueIndex:idx_site_username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	DOB          string `json:"dob"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email" gorm:"uniqueIndex:idx_site_email"`
	Avatar       string `json:"avatar"`
	Role         string `json:"role"`
	Site         string `json:"site" gorm:"uniqueIndex:idx_site_username;uniqueIndex:idx_site_email"`
	TokenVersion int    `json:"token_version" gorm:"default:0"`
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
