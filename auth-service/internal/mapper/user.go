package mapper

import (
	"auth-service/internal/model"
	"auth-service/internal/repository/user"
	api "auth-service/pkg/api/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func ToGRPCResponseUserFromModel(user *model.User) *api.RegisterResponse {
	role := strings.Join(user.Roles, ",")
	return &api.RegisterResponse{
		User: &api.User{
			Uuid:      user.UUID,
			Username:  user.Username,
			Email:     user.Email,
			Roles:     role,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
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
