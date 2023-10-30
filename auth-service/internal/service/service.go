package service

import (
	"auth-service/internal/model"
	"context"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type UserService interface {
	Auth(ctx context.Context, email, password string) (map[string]string, error)
	Create(ctx context.Context, info *model.User) (*model.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (map[string]string, error)
}
