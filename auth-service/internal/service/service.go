package service

import (
	"auth-service/internal/model"
	"context"
)

type UserService interface {
	Auth(ctx context.Context, email, password string) (map[string]string, error)
	Create(ctx context.Context, info *model.User) (map[string]string, error)
	RefreshToken(ctx context.Context, refreshToken string) (map[string]string, error)
}
