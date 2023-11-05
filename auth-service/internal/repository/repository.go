package repository

import (
	"auth-service/internal/model"
	"context"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type UserRepository interface {
	Get(ctx context.Context, userUUID string) *model.User
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) *model.User
	GetByRefreshToken(ctx context.Context, token string) *model.User
	SaveRefreshToken(ctx context.Context, UUID string, session *model.Session) error
}
