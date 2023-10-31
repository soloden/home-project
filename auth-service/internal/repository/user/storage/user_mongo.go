package storage

import (
	"auth-service/internal/config"
	"auth-service/internal/model"
	"auth-service/internal/repository/user"
	"auth-service/internal/service/generation"
	"auth-service/pkg/logger"
	pkgStorage "auth-service/pkg/storage/mongodb"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type mongodbRepository struct {
	mCollection *mongo.Collection
}

func NewMongodbRepository(cfg *config.Config, mongodb *pkgStorage.MongoDB) *mongodbRepository {
	collection := mongodb.Client().Database(cfg.MongoDB.Database).Collection("users")
	return &mongodbRepository{
		mCollection: collection,
	}
}

func (r *mongodbRepository) Create(ctx context.Context, user *user.User, idGen generation.IdGenerator) *user.User {
	log := logger.LoggerFromContext(ctx)
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	if user.UpdatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	uuid := idGen.Generate(ctx)
	if uuid == "" {
		log.Error("generated UUID is empty")
		return nil
	}

	user.UUID = uuid
	bsonData, err := bson.Marshal(user)
	if err != nil {
		log.Error("marshaling for save", zap.Error(err))
		return nil
	}

	_, err = r.mCollection.InsertOne(ctx, bsonData)
	if err != nil {
		log.Error("inserting user", zap.Error(err))
		return nil
	}

	return user
}

func (r *mongodbRepository) Get(ctx context.Context, userUUID string) *user.User {
	var modelUser user.User

	err := r.mCollection.FindOne(
		ctx,
		bson.D{{"uuid", userUUID}},
	).Decode(&modelUser)
	if err != nil {
		logger.LoggerFromContext(ctx).Error("finding user", zap.Error(err))
		return nil
	}

	return &modelUser
}

func (r *mongodbRepository) GetByEmail(ctx context.Context, email string) *user.User {
	var modelUser user.User
	log := logger.LoggerFromContext(ctx)

	err := r.mCollection.FindOne(
		ctx,
		bson.D{{"email", email}},
	).Decode(&modelUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}

		log.Error("mongodb find by email", zap.Error(err))
		return nil
	}

	return &modelUser
}

func (r *mongodbRepository) GetByRefreshToken(ctx context.Context, token string) *user.User {
	var modelUser user.User

	err := r.mCollection.FindOne(
		ctx,
		bson.M{
			"session.refreshToken": token,
			"session.expiresAt":    bson.M{"$gt": time.Now()},
		},
	).Decode(&modelUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}

		logger.LoggerFromContext(ctx).Error("finding user by refresh token", zap.Error(err))
		return nil
	}

	return &modelUser
}

func (r *mongodbRepository) SaveRefreshToken(ctx context.Context, UUID string, session *model.Session) error {
	_, err := r.mCollection.UpdateOne(
		ctx,
		bson.M{"uuid": UUID},
		bson.M{
			"$set": bson.M{"session": session},
		})
	if err != nil {
		logger.LoggerFromContext(ctx).Error("update refresh token", zap.Error(err))
		return err
	}

	return nil
}
