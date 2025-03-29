package user_model

import (
	"encoding/json"
	"fmt"
	"go-auth-service/pkg/shared/dto"
	"io/ioutil"
	"os"
)

type User struct {
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	Name         string   `json:"name"`
	DOB          string   `json:"dob"`
	PhoneNumber  string   `json:"phone_number"`
	Email        string   `json:"email"`
	Avatar       string   `json:"avatar"`
	Role         []string `json:"role"`
	Site         string   `json:"site"`
	TokenVersion int      `json:"token_version"`
}

type UserResponse struct {
	Username    string   `json:"username"`
	Name        string   `json:"name"`
	DOB         string   `json:"dob"`
	PhoneNumber string   `json:"phone"`
	Email       string   `json:"email"`
	Avatar      string   `json:"avatar"`
	Role        []string `json:"role"`
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

func (user User) ToDTO() shared_dto.UserDTO {
	return shared_dto.UserDTO{
		Username:     user.Username,
		Password:     user.Password,
		Name:         user.Name,
		PhoneNumber:  user.PhoneNumber,
		Email:        user.Email,
		Role:         user.Role,
		Site:         user.Site,
		TokenVersion: user.TokenVersion,
	}
}

func (user User) HasRole(role string) bool {
	for _, r := range user.Role {
		if r == role {
			return true
		}
	}
	return false
}

func (user User) HasAnyRole(roles []string) bool {
	for _, requiredRole := range roles {
		for _, userRole := range user.Role {
			if userRole == requiredRole {
				return true
			}
		}
	}
	return false
}

var UserList []User

func (user User) IncreaseTokenVersion() {
	user.TokenVersion++
}

func LoadUsersFromFile(filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist")
	}

	// Read the file contents
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into the UserList slice
	err = json.Unmarshal(fileData, &UserList)
	if err != nil {
		return err
	}

	return nil
}
