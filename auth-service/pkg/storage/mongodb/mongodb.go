package mongodb

import (
	"auth-service/internal/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client mongo.Client
}

var mongoInstance *MongoDB

func NewClient(ctx context.Context) (*MongoDB, error) {
	if mongoInstance != nil {
		return mongoInstance, nil
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("loading config: %s", err)
	}
	clientOptions := options.Client().ApplyURI(cfg.MongoDB.URL)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("initialization mongodb connect: %s", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("pinging mongodb: %s", err)
	}

	mongoInstance = &MongoDB{
		client: *client,
	}

	return mongoInstance, nil
}

func (mongodb *MongoDB) Client() *mongo.Client {
	return &mongodb.client
}

func (mongodb MongoDB) CloseClient(ctx context.Context) error {
	err := mongodb.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("close mongodb connection: %s", err)
	}

	return nil
}
