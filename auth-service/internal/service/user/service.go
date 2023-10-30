package user

import (
	"auth-service/internal/mapper"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/internal/service/auth"
	"auth-service/pkg/logger"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

func (us *userService) Create(ctx context.Context, user *model.User) (*model.User, error) {
	log := logger.LoggerFromContext(ctx)
	userUUID, err := uuid.NewUUID()
	if err != nil {
		log.Error("creating UUID problem", zap.Error(err))
		return nil, status.Error(codes.Internal, "something went wrong, please try later")
	}

	checkUser := us.userRepository.GetByEmail(ctx, user.Email)
	if checkUser != nil {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	user.UUID = userUUID.String()
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	if err != nil {
		log.Error("hashing user password", zap.Error(err))
		return nil, status.Error(codes.Internal, "something went wrong, please try later")
	}

	repUser := mapper.ToUserFromModel(user)
	repUser.PasswordHash = hashPassword

	savedUser := us.userRepository.Create(ctx, repUser)
	if savedUser == nil {
		return nil, status.Error(codes.Internal, "something went wrong, please try later")
	}

	return mapper.ToUserFromRepository(savedUser), nil
}

func (us *userService) Auth(ctx context.Context, email, password string) (map[string]string, error) {
	user := us.userRepository.GetByEmail(ctx, email)
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "wrong password")
	}

	res, err := us.generateTokens(ctx, mapper.ToUserFromRepository(user))
	if err != nil {
		return nil, status.Error(codes.Internal, "something went wrong, please try later")
	}

	return res, nil
}

func (us *userService) RefreshToken(ctx context.Context, refreshToken string) (map[string]string, error) {
	user := us.userRepository.GetByRefreshToken(ctx, refreshToken)
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	res, err := us.generateTokens(ctx, mapper.ToUserFromRepository(user))
	if err != nil {
		return nil, status.Error(codes.Internal, "something went wrong, please try later")
	}

	return res, nil
}

func (us *userService) generateTokens(ctx context.Context, user *model.User) (map[string]string, error) {
	log := logger.LoggerFromContext(ctx)
	jwtToken, err := auth.GenerateJWTToken(auth.ToUserFromModel(user))
	if err != nil {
		log.Error("generate JWT token", zap.Error(err))
		return nil, fmt.Errorf("auth user: %s", err)
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		log.Error("generate refresh token", zap.Error(err))
		return nil, fmt.Errorf("generation tokens: %s", err)
	}

	session := &model.Session{
		refreshToken,
		time.Now().Add(time.Hour * 24),
	}
	err = us.userRepository.SaveRefreshToken(ctx, user.UUID, session)
	if err != nil {
		return nil, err
	}

	tokens := make(map[string]string)
	tokens["token"] = jwtToken
	tokens["refreshToken"] = refreshToken

	return tokens, nil
}
