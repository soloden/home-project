package generation

import (
	"auth-service/pkg/logger"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IdGenerator interface {
	Generate(ctx context.Context) string
}

//go:generate mockgen -source=uuid.go -destination=mocks/mock.go
type UUIDGeneration struct{}

func (ug UUIDGeneration) Generate(ctx context.Context) string {
	log := logger.LoggerFromContext(ctx)
	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Error("generation uuid failure", zap.Error(err))
		return ""
	}

	return uuid.String()
}
