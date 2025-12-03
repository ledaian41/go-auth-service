package auth_model

import "go-auth-service/internal/shared/interface"

type AuthService struct {
	userService shared_interface.UserService
	secretKey   string
}
