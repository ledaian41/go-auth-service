package auth

import (
	"errors"
	"go-auth-service/config"
	"go-auth-service/internal/shared"
	"strings"
)

type Service struct {
	userService  shared.UserService
	tokenService shared.TokenService
	redisClient  *config.RedisClient
}

func NewAuthService(redisClient *config.RedisClient, userService shared.UserService, tokenService shared.TokenService) *Service {
	return &Service{
		redisClient:  redisClient,
		userService:  userService,
		tokenService: tokenService,
	}
}

func (s *Service) CheckValidUser(username, password, siteId string) (*shared.UserDTO, error) {
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

	if !shared.CheckHashPassword(password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *Service) CreateNewAccount(account *shared.RegisterRequestDTO) (*shared.UserDTO, error) {
	hashedPassword, err := shared.HashPassword(account.Password)
	if err != nil {
		return nil, err
	}

	newUser := shared.UserDTO{
		Username:    account.Username,
		Password:    hashedPassword,
		PhoneNumber: account.PhoneNumber,
		Email:       account.Email,
		Role:        "user",
	}
	return s.userService.CreateNewUser(&newUser)
}

func (s *Service) FindUserByUsername(username, siteId string) (*shared.UserDTO, error) {
	return s.userService.FindUserByUsername(username, siteId)
}

func (s *Service) CheckAdminRole(role []interface{}) bool {
	for _, r := range role {
		if r == "admin" {
			return true
		}
	}
	return false
}
