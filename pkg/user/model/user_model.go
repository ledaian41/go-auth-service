package user_model

import "errors"

type User struct {
	Username    string
	Password    string
	Name        string
	DOB         string
	PhoneNumber string
	Email       string
	Avatar      string
	Role        []string
	Site        string
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

func (user User) Validate() (bool, error) {
	return true, nil
}

var UserList = []User{
	{
		Username:    "admin",
		Password:    "$2a$10$efVRe6.fsZL41t.3Nxp61OaCqS40pdUyP7LOmxccGnceisga6iovG",
		Name:        "Harry",
		DOB:         "11/07/2024",
		PhoneNumber: "0703940225",
		Email:       "harry@gmail.com",
		Avatar:      "",
		Role:        []string{"admin", "user"},
		Site:        "lexis",
	},
	{
		Username:    "An",
		Password:    "$2a$10$efVRe6.fsZL41t.3Nxp61OaCqS40pdUyP7LOmxccGnceisga6iovG",
		Name:        "Le Dai An",
		DOB:         "04/01/1995",
		PhoneNumber: "0703940225",
		Email:       "ledaian41@gmail.com",
		Avatar:      "",
		Role:        []string{"manager", "user"},
		Site:        "lexis",
	}}

func GetById(userId string) (*User, error) {
	for _, user := range UserList {
		if user.Username == userId {
			userCopy := user
			return &userCopy, nil
		}
	}
	return nil, errors.New("user not found")
}
