package auth_v1

import (
	"auth-service/internal/mapper"
	"auth-service/internal/service"
	api "auth-service/pkg/api/auth_v1"
	"auth-service/pkg/logger"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	userService service.UserService
	api.UnimplementedAuthServer
}

func NewAuthServer(service service.UserService) *AuthServer {
	return &AuthServer{
		userService: service,
	}
}

func (as *AuthServer) Login(ctx context.Context, req *api.AuthRequest) (*api.Tokens, error) {
	logger.ContextWithLogger(ctx, zap.L())
	email := req.GetEmail()
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	password := req.GetPassword()
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	res, err := as.userService.Auth(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &api.Tokens{
		Token:        res["token"],
		RefreshToken: res["refreshToken"],
	}, nil
}

func (as *AuthServer) Register(ctx context.Context, req *api.RegisterRequest) (*api.Tokens, error) {
	logger.ContextWithLogger(ctx, zap.L())
	email := req.GetUser().GetEmail()
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	password := req.GetUser().GetPassword()
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	res, err := as.userService.Create(ctx, mapper.ToUserFromGRPSRequest(req.GetUser()))
	if err != nil {
		return nil, err
	}

	response := &api.Tokens{Token: res["token"], RefreshToken: res["refreshToken"]}
	return response, nil
}

func (as *AuthServer) RefreshToken(ctx context.Context, req *api.RefreshTokenRequest) (*api.Tokens, error) {
	logger.ContextWithLogger(ctx, zap.L())
	token := req.GetRefreshToken()
	if token == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	res, err := as.userService.RefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &api.Tokens{
		Token:        res["token"],
		RefreshToken: res["refreshToken"],
	}, nil
}
