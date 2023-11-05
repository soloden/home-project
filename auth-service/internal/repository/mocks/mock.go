// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=mocks/mock.go
//
// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	model "auth-service/internal/model"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUserRepositoryMockRecorder) Create(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserRepository)(nil).Create), ctx, user)
}

// Get mocks base method.
func (m *MockUserRepository) Get(ctx context.Context, userUUID string) *model.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, userUUID)
	ret0, _ := ret[0].(*model.User)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockUserRepositoryMockRecorder) Get(ctx, userUUID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserRepository)(nil).Get), ctx, userUUID)
}

// GetByEmail mocks base method.
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) *model.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", ctx, email)
	ret0, _ := ret[0].(*model.User)
	return ret0
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockUserRepositoryMockRecorder) GetByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockUserRepository)(nil).GetByEmail), ctx, email)
}

// GetByRefreshToken mocks base method.
func (m *MockUserRepository) GetByRefreshToken(ctx context.Context, token string) *model.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByRefreshToken", ctx, token)
	ret0, _ := ret[0].(*model.User)
	return ret0
}

// GetByRefreshToken indicates an expected call of GetByRefreshToken.
func (mr *MockUserRepositoryMockRecorder) GetByRefreshToken(ctx, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByRefreshToken", reflect.TypeOf((*MockUserRepository)(nil).GetByRefreshToken), ctx, token)
}

// SaveRefreshToken mocks base method.
func (m *MockUserRepository) SaveRefreshToken(ctx context.Context, UUID string, session *model.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveRefreshToken", ctx, UUID, session)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveRefreshToken indicates an expected call of SaveRefreshToken.
func (mr *MockUserRepositoryMockRecorder) SaveRefreshToken(ctx, UUID, session any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveRefreshToken", reflect.TypeOf((*MockUserRepository)(nil).SaveRefreshToken), ctx, UUID, session)
}
