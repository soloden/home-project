package app

import (
	"auth-service/internal/api/auth_v1"
	"auth-service/internal/config"
	"auth-service/internal/repository"
	"auth-service/internal/repository/user/storage"
	"auth-service/internal/service"
	"auth-service/internal/service/user"
	api "auth-service/pkg/api/auth_v1"
	"auth-service/pkg/storage/mongodb"
	"context"
)

type serviceProvider struct {
	userRepository repository.UserRepository
	userService    service.UserService
	authServer     api.AuthServer
	mongodb        *mongodb.MongoDB
}

func (sp *serviceProvider) InitUserRepository() repository.UserRepository {
	if sp.userRepository == nil {
		cfg, _ := config.LoadConfig()
		switch cfg.App.StorageType {
		case "mongodb":
			sp.userRepository = storage.NewMongodbRepository(cfg, sp.InitMongoDBConnection(context.TODO()))
		default:
			sp.userRepository = storage.NewMemoryRepository()
		}
	}

	return sp.userRepository
}

func (sp *serviceProvider) InitUserService() service.UserService {
	if sp.userService == nil {
		sp.userService = user.NewUserService(sp.InitUserRepository())
	}

	return sp.userService
}

func (sp *serviceProvider) InitAuthServer() api.AuthServer {
	if sp.authServer == nil {
		sp.authServer = auth_v1.NewAuthServer(sp.InitUserService())
	}

	return sp.authServer
}

func (sp *serviceProvider) InitMongoDBConnection(ctx context.Context) *mongodb.MongoDB {
	if sp.mongodb == nil {
		client, err := mongodb.NewClient(ctx)
		if err != nil {
			panic(err)
		}

		sp.mongodb = client
	}

	return sp.mongodb
}
