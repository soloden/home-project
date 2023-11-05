package user

import (
	"auth-service/internal/model"
	repositoryMock "auth-service/internal/repository/mocks"
	generationMock "auth-service/internal/service/generation/mocks"
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestUser_CreateSuccess(t *testing.T) {
	password := "test123"
	c := gomock.NewController(t)
	defer c.Finish()
	uuidGen := generationMock.NewMockIdGenerator(c)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		t.Fatalf("hashing user password %v", err)
	}

	payload := &model.User{
		Username: "test",
		Email:    "test@test.org",
		Password: password,
		Roles:    []string{"admin", "test"},
	}

	expected := &model.User{
		UUID:         uuidGen.Generate(context.TODO()),
		Email:        "test@test.org",
		Username:     "test",
		Password:     password,
		PasswordHash: hash,
		Roles:        []string{"admin", "test"},
	}

	repo := repositoryMock.NewMockUserRepository(c)
	repo.EXPECT().GetByEmail(context.TODO(), payload.Email).Return(nil).Times(1)
	repo.EXPECT().Create(context.TODO(), payload).Return(nil).Times(1)
	service := NewUserService(repo)
	res, err := service.Create(context.TODO(), payload, uuidGen)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}
