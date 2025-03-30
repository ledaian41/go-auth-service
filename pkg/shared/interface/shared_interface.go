package shared_interface

import (
	"github.com/golang-jwt/jwt/v5"
	"go-auth-service/pkg/shared/dto"
)

type SiteService interface {
	CheckSiteExists(siteId string) *shared_dto.SiteDTO
}

type AuthService interface {
	CreateNewAccount(account *shared_dto.RegisterRequestDTO) (*shared_dto.UserDTO, error)
	GenerateAccessToken(siteSecretKey, sessionId string, user *shared_dto.UserDTO) (string, error)
	ValidateAccessToken(site *shared_dto.SiteDTO, tokenStr string) (jwt.MapClaims, error)
	GenerateRefreshToken(user *shared_dto.UserDTO) (string, error)
	ValidateRefreshToken(tokenStr string) (jwt.MapClaims, error)
	CheckValidUser(username, password, siteId string) (*shared_dto.UserDTO, error)
	FindUserByUsername(username, siteId string) (*shared_dto.UserDTO, error)
	CheckAdminRole(role []interface{}) bool
	RevokeSessionId(sessionId string)
}

type UserService interface {
	CreateNewUser(user *shared_dto.UserDTO) (*shared_dto.UserDTO, error)
	FindUserByUsername(username, siteId string) (*shared_dto.UserDTO, error)
	FindUsersBySite(siteId string) *[]shared_dto.UserDTO
}

type TokenService interface {
	ValidateRefreshToken(refreshToken string) string
	StoreRefreshToken(username, refreshToken string) string
	RevokeRefreshToken(refreshToken string) string
}
