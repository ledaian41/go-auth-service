package user_service

import (
	"auth/internal/user/model"
	"errors"
)

func CreateNewUser(user *user_model.User) (*user_model.User, error) {
	if IsUsernameExist(user.Username) {
		return nil, errors.New("username exist")
	}

	user_model.UserList = append(user_model.UserList, *user)
	return user, nil
}

func IsUsernameExist(username string) bool {
	for _, user := range user_model.UserList {
		if user.Username == username {
			return true
		}
	}
	return false
}
