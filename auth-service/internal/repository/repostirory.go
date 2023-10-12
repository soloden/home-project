package repository

import (
	"auth-service/internal/model"
	"auth-service/internal/repository/user"
	"context"
)

type UserRepository interface {
	Get(ctx context.Context, userUUID string) *user.User
	Create(ctx context.Context, user *user.User) *user.User
	GetByEmail(ctx context.Context, email string) *user.User
	GetByRefreshToken(ctx context.Context, token string) *user.User
	SaveRefreshToken(ctx context.Context, UUID string, session *model.Session) error
}
