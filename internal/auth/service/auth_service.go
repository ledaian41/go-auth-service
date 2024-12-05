package auth_service

import (
	"auth/internal/auth/model"
	"auth/internal/user/model"
	"auth/internal/user/service"
	"errors"
	"strings"
)

func CheckValidUser(username string, password string) (*user_model.User, error) {
	if len(strings.Trim(username, " ")) == 0 {
		return nil, errors.New("invalid username or password")
	}

	if len(strings.Trim(password, " ")) == 0 {
		return nil, errors.New("invalid username or password")
	}

	user, err := user_model.GetById(username)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func GetSessionUser(username string) (*user_model.User, error) {
	if len(strings.Trim(username, " ")) == 0 {
		return nil, errors.New("user not found")
	}

	user, err := user_model.GetById(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateNewAccount(account *auth_model.RegisterAccount) (*user_model.User, error) {
	newUser := user_model.User{
		Username:    account.Username,
		Password:    account.Password,
		PhoneNumber: account.PhoneNumber,
		Email:       account.Email,
	}
	user, err := user_service.CreateNewUser(&newUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}
