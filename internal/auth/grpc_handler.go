package auth

import (
	"context"
	"errors"
	"go-auth-service/internal/shared"
	auth "go-auth-service/proto"
)

type GrpcHandler struct {
	auth.UnimplementedAuthServer
	siteService  shared.SiteService
	authService  shared.AuthService
	tokenService shared.TokenService
}

func New(siteService shared.SiteService, authService shared.AuthService, tokenService shared.TokenService) *GrpcHandler {
	return &GrpcHandler{siteService: siteService, authService: authService, tokenService: tokenService}
}

func (handler *GrpcHandler) Ping(_ context.Context, req *auth.HelloRequest) (*auth.HelloResponse, error) {
	return &auth.HelloResponse{Message: "Hello, " + req.Name}, nil
}

func (handler *GrpcHandler) Login(_ context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	siteId := req.Site
	if siteId == "" {
		siteId = "app"
	}

	user, err := handler.authService.CheckValidUser(req.Username, req.Password, siteId)
	if err != nil {
		return nil, err
	}

	sessionId := shared.RandomID()
	refreshToken, err := handler.authService.GenerateRefreshToken(user, sessionId)
	if err != nil {
		return nil, err
	}

	_, err = handler.tokenService.StoreRefreshToken(user.Username, sessionId)
	if err != nil {
		return nil, errors.New("error saving refresh token")
	}

	site := handler.siteService.CheckSiteExists(siteId)
	if site == nil {
		return nil, errors.New("site not found")
	}

	accessToken, err := handler.authService.GenerateAccessToken(site.SecretKey, sessionId, user)
	return &auth.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (handler *GrpcHandler) RefreshToken(_ context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, errors.New("invalid refresh token")
	}

	siteId := req.Site
	if siteId == "" {
		siteId = "app"
	}

	site := handler.siteService.CheckSiteExists(siteId)
	if site == nil {
		return nil, errors.New("site not found")
	}

	newAccessToken, newRefreshToken, err := handler.authService.RotateRefreshToken(site, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &auth.RefreshTokenResponse{AccessToken: newAccessToken, RefreshToken: newRefreshToken}, nil
}

func (handler *GrpcHandler) JWT(_ context.Context, req *auth.JwtRequest) (*auth.JwtResponse, error) {
	siteId := req.Site
	if siteId == "" {
		siteId = "app"
	}

	site := handler.siteService.CheckSiteExists(siteId)
	if site == nil {
		return nil, errors.New("site not found")
	}

	if req.AccessToken == "" {
		return nil, errors.New("invalid access token")
	}

	claims, err := handler.authService.ValidateAccessToken(site, req.AccessToken)
	if err != nil {
		return nil, err
	}

	return &auth.JwtResponse{
		Username: claims["user"].(string),
		Role:     claims["role"].(string),
		Name:     claims["name"].(string),
		Email:    claims["email"].(string),
		Phone:    claims["phone"].(string),
	}, nil
}

func (handler *GrpcHandler) Logout(_ context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	claims, err := handler.authService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return &auth.LogoutResponse{Status: false}, nil
	}

	sessionId := claims["jti"].(string)
	handler.tokenService.RevokeRefreshToken(sessionId)
	handler.authService.RevokeSessionId(sessionId)

	return &auth.LogoutResponse{
		Status: true,
	}, nil
}
