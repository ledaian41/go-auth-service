package auth_service

import (
	"errors"
	"go-auth-service/config"
	"go-auth-service/pkg/shared/dto"
	"go-auth-service/pkg/shared/interface"
	"go-auth-service/pkg/shared/utils"
	"strings"
	"time"
)

var (
	accessExpireTime  = time.Minute * 15   // 15 minutes
	refreshExpireTime = time.Hour * 24 * 7 // 1 week
)

type AuthService struct {
	userService shared_interface.UserServiceInterface
	secretKey   string
	redisClient *config.RedisClient
}

func NewAuthService(userService shared_interface.UserServiceInterface, redisClient *config.RedisClient) *AuthService {
	return &AuthService{
		userService: userService,
		secretKey:   config.Env.SecretKey,
		redisClient: redisClient,
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
		Role:        []string{"user"},
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
