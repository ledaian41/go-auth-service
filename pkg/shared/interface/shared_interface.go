package shared_interface

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-auth-service/pkg/shared/dto"
)

type SiteServiceInterface interface {
	CheckSiteExists(siteId string) *shared_dto.SiteDTO
}

type AuthServiceInterface interface {
	CreateNewAccount(account *shared_dto.RegisterRequestDTO) (*shared_dto.UserDTO, error)
	GenerateAccessToken(c *gin.Context, user *shared_dto.UserDTO) (string, error)
	ValidateAccessToken(c *gin.Context, tokenStr string) (jwt.MapClaims, error)
	GenerateRefreshToken(user *shared_dto.UserDTO) (string, error)
	ValidateRefreshToken(tokenStr string) (jwt.MapClaims, error)
	CheckValidUser(username, password, siteId string) (*shared_dto.UserDTO, error)
	FindUserByUsername(username, siteId string) (*shared_dto.UserDTO, error)
	CheckAdminRole(role []interface{}) bool
	RevokeUserSession(username string)
}

type UserServiceInterface interface {
	CreateNewUser(user *shared_dto.UserDTO) (*shared_dto.UserDTO, error)
	FindUserByUsername(username, siteId string) (*shared_dto.UserDTO, error)
	FindUsersBySite(siteId string) *[]shared_dto.UserDTO
	IncrementTokenVersion(username string)
}
