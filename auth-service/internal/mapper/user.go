package mapper

import (
	"auth-service/internal/model"
	"auth-service/internal/repository/user"
	api "auth-service/pkg/api/auth_v1"
	"strings"
)

func ToUserFromGRPSRequest(info *api.User) *model.User {
	roleList := strings.Split(info.GetRoles(), ",")
	return &model.User{
		Email:    info.GetEmail(),
		Username: info.GetUsername(),
		Password: info.GetPassword(),
		Roles:    roleList,
	}
}

func ToUserFromRepository(user *user.User) *model.User {
	return &model.User{
		UUID:      user.UUID,
		Email:     user.Email,
		Username:  user.Username,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserFromModel(mUser *model.User) *user.User {
	return &user.User{
		UUID:      mUser.UUID,
		Email:     mUser.Email,
		Username:  mUser.Username,
		Roles:     mUser.Roles,
		CreatedAt: mUser.CreatedAt,
		UpdatedAt: mUser.UpdatedAt,
	}
}
