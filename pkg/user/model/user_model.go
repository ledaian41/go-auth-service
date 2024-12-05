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
	Role        string
	Site        string
}

type UserResponse struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	DOB         string `json:"dob"`
	PhoneNumber string `json:"phone"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
}

func (user User) Response() UserResponse {
	return UserResponse{
		Username:    user.Username,
		Name:        user.Name,
		DOB:         user.DOB,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Avatar:      user.Avatar,
	}
}

func (user User) Validate() (bool, error) {
	return true, nil
}

var UserList = []User{
	{
		Username:    "admin",
		Password:    "123456",
		Name:        "Harry",
		DOB:         "11/07/2024",
		PhoneNumber: "0703940225",
		Email:       "harry@gmail.com",
		Avatar:      "",
		Role:        "admin",
		Site:        "lexis",
	},
	{
		Username:    "An",
		Password:    "123456",
		Name:        "Le Dai An",
		DOB:         "04/01/1995",
		PhoneNumber: "0703940225",
		Email:       "ledaian41@gmail.com",
		Avatar:      "",
		Role:        "manager",
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
