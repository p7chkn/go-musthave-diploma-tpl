// Code generated by mockery v2.9.4. DO NOT EDIT.

package handlers

import (
	context "context"

	models "github.com/p7chkn/go-musthave-diploma-tpl/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MockRepositoryInterface is an autogenerated mock type for the RepositoryInterface type
type MockRepositoryInterface struct {
	mock.Mock
}

// CheckPassword provides a mock function with given fields: ctx, user
func (_m *MockRepositoryInterface) CheckPassword(ctx context.Context, user models.User) (models.User, error) {
	ret := _m.Called(ctx, user)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(context.Context, models.User) models.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateOrder provides a mock function with given fields: ctx, order
func (_m *MockRepositoryInterface) CreateOrder(ctx context.Context, order models.Order) error {
	ret := _m.Called(ctx, order)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *MockRepositoryInterface) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	ret := _m.Called(ctx, user)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(context.Context, models.User) *models.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrders provides a mock function with given fields: ctx, userID
func (_m *MockRepositoryInterface) GetOrders(ctx context.Context, userID string) ([]models.Order, error) {
	ret := _m.Called(ctx, userID)

	var r0 []models.Order
	if rf, ok := ret.Get(0).(func(context.Context, string) []models.Order); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Order)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Ping provides a mock function with given fields: ctx
func (_m *MockRepositoryInterface) Ping(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
