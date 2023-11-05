package storage

import (
	"auth-service/internal/model"
	"context"
	"fmt"
	"sync"
)

type memoryRepository struct {
	data  map[string]*model.User
	mutex sync.RWMutex
}

func NewMemoryRepository() *memoryRepository {
	return &memoryRepository{
		data: make(map[string]*model.User),
	}
}

func (r *memoryRepository) Create(ctx context.Context, mUser *model.User) error {
	data := mUser

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.data[data.UUID] = data
	fmt.Println(r.data)
	return nil
}

func (r *memoryRepository) Get(_ context.Context, userUUID string) *model.User {
	r.mutex.RLock()
	defer r.mutex.Unlock()

	result, ok := r.data[userUUID]
	if ok {
		return result
	}

	return nil
}

func (r *memoryRepository) GetByEmail(_ context.Context, _ string) *model.User {
	return nil
}

func (r *memoryRepository) GetByRefreshToken(_ context.Context, _ string) *model.User {
	return nil
}

func (r *memoryRepository) SaveRefreshToken(_ context.Context, _ string, _ *model.Session) error {
	return nil
}
