package service

import (
	"auth-service/internal/model"
	"auth-service/internal/service/generation"
	"context"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type UserService interface {
	Auth(ctx context.Context, email, password string) (map[string]string, error)
	Create(ctx context.Context, info *model.User, generator generation.IdGenerator) (*model.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (map[string]string, error)
}
