package auth_model

import "go-auth-service/pkg/shared/interface"

type AuthService struct {
	userService shared_interface.UserServiceInterface
	secretKey   string
}
