package shared

import (
	"github.com/golang-jwt/jwt/v5"
)

type SiteService interface {
	CheckSiteExists(siteId string) *SiteDTO
}

type AuthService interface {
	CreateNewAccount(account *RegisterRequestDTO) (*UserDTO, error)
	GenerateAccessToken(siteSecretKey, sessionId string, user *UserDTO) (string, error)
	ValidateAccessToken(site *SiteDTO, tokenStr string) (jwt.MapClaims, error)
	GenerateRefreshToken(user *UserDTO, sessionId string) (string, error)
	ValidateRefreshToken(tokenStr string) (jwt.MapClaims, error)
	CheckValidUser(username, password, siteId string) (*UserDTO, error)
	CheckAdminRole(role []interface{}) bool
	RevokeSessionId(sessionId string)
	RotateRefreshToken(site *SiteDTO, oldRefreshToken string) (newAccessToken, newRefreshToken string, err error)
}

type UserService interface {
	CreateNewUser(user *UserDTO) (*UserDTO, error)
	FindUserByUsername(username, siteId string) (*UserDTO, error)
	FindUsersBySite(siteId string) *[]UserDTO
}

type TokenService interface {
	ValidateRefreshToken(id string) string
	StoreRefreshToken(username, id string) (string, error)
	RevokeRefreshToken(id string) string
}
