package auth_service

import (
	"errors"
	"go-auth-service/config"
	"go-auth-service/internal/shared/dto"
	"go-auth-service/internal/shared/interface"
	"go-auth-service/internal/shared/utils"
	"strings"
)

type AuthService struct {
	userService shared_interface.UserService
	redisClient *config.RedisClient
}

func NewAuthService(redisClient *config.RedisClient, userService shared_interface.UserService) *AuthService {
	return &AuthService{
		redisClient: redisClient,
		userService: userService,
	}
}

func (s *AuthService) CheckValidUser(username, password, siteId string) (*shared_dto.UserDTO, error) {
	if len(strings.Trim(username, " ")) == 0 {
		return nil, errors.New("invalid username or password")
	}

	if len(strings.Trim(password, " ")) == 0 {
		return nil, errors.New("invalid username or password")
	}

	user, err := s.FindUserByUsername(username, siteId)
	if err != nil {
		return nil, err
	}

	if !shared_utils.CheckHashPassword(password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *AuthService) CreateNewAccount(account *shared_dto.RegisterRequestDTO) (*shared_dto.UserDTO, error) {
	hashedPassword, err := shared_utils.HashPassword(account.Password)
	if err != nil {
		return nil, err
	}

	newUser := shared_dto.UserDTO{
		Username:    account.Username,
		Password:    hashedPassword,
		PhoneNumber: account.PhoneNumber,
		Email:       account.Email,
		Role:        "user",
	}
	return s.userService.CreateNewUser(&newUser)
}

func (s *AuthService) FindUserByUsername(username, siteId string) (*shared_dto.UserDTO, error) {
	return s.userService.FindUserByUsername(username, siteId)
}

func (s *AuthService) CheckAdminRole(role []interface{}) bool {
	for _, r := range role {
		if r == "admin" {
			return true
		}
	}
	return false
}
