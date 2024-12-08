package user_service

import (
	"errors"
	shared_dto "go-auth-service/pkg/shared/dto"
	shared_utils "go-auth-service/pkg/shared/utils"
	user_model "go-auth-service/pkg/user/model"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) FindUserByUsername(username string) (*shared_dto.UserDTO, error) {
	for _, user := range user_model.UserList {
		if user.Username == username {
			userCopy := user.ToDTO()
			return &userCopy, nil
		}
	}
	return nil, errors.New("user not found")
}

func (s *UserService) CreateNewUser(user *shared_dto.UserDTO) (*shared_dto.UserDTO, error) {
	if IsUsernameExist(user.Username) {
		return nil, errors.New("username exist")
	}

	newUser := user_model.User{
		Username:    user.Username,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        user.Role,
		Site:        user.Site,
	}
	user_model.UserList = append(user_model.UserList, newUser)
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

func (s *UserService) FindUsersBySite(siteId string) *[]shared_dto.UserDTO {
	filteredUsers := shared_utils.Filter(user_model.UserList, func(user user_model.User) bool {
		return user.Site == siteId
	})
	result := shared_utils.Map(filteredUsers, func(user user_model.User) shared_dto.UserDTO {
		return shared_dto.UserDTO{
			Username:    user.Username,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
			Role:        user.Role,
		}
	})
	return &result
}
