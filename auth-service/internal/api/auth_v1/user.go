package auth_v1

import (
	"auth-service/internal/model"
	api "auth-service/pkg/api/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

func (as *AuthServer) RegisterRequestToDTO(req *api.RegisterRequest) model.User {
	roleList := strings.Split(req.GetUser().GetRoles(), ",")
	return model.User{
		Email:    req.GetUser().GetEmail(),
		Username: req.GetUser().GetUsername(),
		Password: req.GetUser().GetPassword(),
		Roles:    roleList,
	}
}

func (as *AuthServer) RegisterResponseFromDTO(user *model.User) api.RegisterResponse {
	roles := strings.Join(user.Roles, ",")
	return api.RegisterResponse{
		User: &api.User{
			Uuid:      user.UUID,
			Username:  user.Username,
			Email:     user.Email,
			Roles:     roles,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}
}

func (as *AuthServer) ValidateRegisterRequest(req *api.RegisterRequest) error {
	if req.GetUser().GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	password := req.GetUser().GetPassword()
	if password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func (as *AuthServer) ValidateAuthRequest(req *api.AuthRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	password := req.GetPassword()
	if password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}
