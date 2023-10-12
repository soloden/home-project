package storage

import (
	"auth-service/internal/model"
	"auth-service/internal/repository/user"
	"context"
	"fmt"
	"sync"
)

type memoryRepository struct {
	data  map[string]*user.User
	mutex sync.RWMutex
}

func NewMemoryRepository() *memoryRepository {
	return &memoryRepository{
		data: make(map[string]*user.User),
	}
}

func (r *memoryRepository) Create(_ context.Context, mUser *user.User) *user.User {
	data := mUser

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.data[data.UUID] = data
	fmt.Println(r.data)
	return data
}

func (r *memoryRepository) Get(_ context.Context, userUUID string) *user.User {
	r.mutex.RLock()
	defer r.mutex.Unlock()

	result, ok := r.data[userUUID]
	if ok {
		return result
	}

	return nil
}

func (r *memoryRepository) GetByEmail(_ context.Context, _ string) *user.User {
	return nil
}

func (r *memoryRepository) GetByRefreshToken(_ context.Context, _ string) *user.User {
	return nil
}
func (r *memoryRepository) SaveRefreshToken(_ context.Context, _ string, _ *model.Session) error {
	return nil
}
