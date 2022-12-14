// Code generated by mockery v2.15.0. DO NOT EDIT.

package user

import (
	context "context"

	config "gitlab.com/krespix/gamification-api/internal/core/config"

	mock "github.com/stretchr/testify/mock"

	models "gitlab.com/krespix/gamification-api/internal/models"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, user
func (_m *MockService) Create(ctx context.Context, user *models.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, id
func (_m *MockService) Get(ctx context.Context, id int64) (*models.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(context.Context, int64) *models.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InitSuperAdmin provides a mock function with given fields: ctx, admin
func (_m *MockService) InitSuperAdmin(ctx context.Context, admin config.SuperAdmin) error {
	ret := _m.Called(ctx, admin)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, config.SuperAdmin) error); ok {
		r0 = rf(ctx, admin)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// List provides a mock function with given fields: ctx
func (_m *MockService) List(ctx context.Context) ([]*models.User, error) {
	ret := _m.Called(ctx)

	var r0 []*models.User
	if rf, ok := ret.Get(0).(func(context.Context) []*models.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockService creates a new instance of MockService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockService(t mockConstructorTestingTNewMockService) *MockService {
	mock := &MockService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
