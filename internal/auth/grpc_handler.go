package auth

import (
	"context"
	"errors"
	shared_interface "go-auth-service/internal/shared/interface"
	"go-auth-service/proto"
	"log"
)

type GrpcHandler struct {
	auth.UnimplementedAuthServer
	siteService  shared_interface.SiteService
	authService  shared_interface.AuthService
	tokenService shared_interface.TokenService
}

func New(siteService shared_interface.SiteService, authService shared_interface.AuthService, tokenService shared_interface.TokenService) *GrpcHandler {
	return &GrpcHandler{siteService: siteService, authService: authService, tokenService: tokenService}
}

func (handler *GrpcHandler) Ping(_ context.Context, req *auth.HelloRequest) (*auth.HelloResponse, error) {
	log.Printf("Received: %s", req.Name)
	return &auth.HelloResponse{Message: "Hello, " + req.Name}, nil
}

func (handler *GrpcHandler) Login(_ context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	log.Printf("Received: %s - %s", req.Username, req.Password)
	siteId := req.Site
	if siteId == "" {
		siteId = "app"
	}

	user, err := handler.authService.CheckValidUser(req.Username, req.Password, siteId)
	if err != nil {
		return nil, err
	}

	refreshToken, err := handler.authService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	sessionId := handler.tokenService.StoreRefreshToken(user.Username, refreshToken)
	if len(sessionId) == 0 {
		return nil, errors.New("error saving refresh token")
	}

	site := handler.siteService.CheckSiteExists(siteId)
	accessToken, err := handler.authService.GenerateAccessToken(site.SecretKey, sessionId, user)
	return &auth.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
