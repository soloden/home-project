package app

import (
	"auth-service/internal/config"
	api "auth-service/pkg/api/auth_v1"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	logger          *zap.Logger
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialization deps: %s", err)
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return fmt.Errorf("init deps: %s", err)
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	_, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("init a.Config: %s", err)
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = &serviceProvider{}
	return nil
}

func (a *App) initGRPCServer(_ context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	api.RegisterAuthServer(a.grpcServer, a.serviceProvider.InitAuthServer())

	return nil
}

func (a *App) Run() error {
	err := a.StartGRPCServer()
	if err != nil {
		return fmt.Errorf("running GRPC server: %v", err)
	}

	return nil
}

func (a *App) StartGRPCServer() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("loading config: %v", err)
	}

	listener, err := net.Listen("tcp", cfg.HttpServer.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	fmt.Printf("server listening at %v\n", listener.Addr())
	if err = a.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("serve GRPS server: %v", err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	err := a.serviceProvider.mongodb.CloseClient(ctx)
	if err != nil {
		return fmt.Errorf("close mongodb connection: %s", err)
	}

	return nil
}
