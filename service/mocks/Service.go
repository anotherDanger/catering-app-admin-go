// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	domain "catering-admin-go/domain"
	context "context"

	mock "github.com/stretchr/testify/mock"

	web "catering-admin-go/web"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// AddProduct provides a mock function with given fields: ctx, request
func (_m *Service) AddProduct(ctx context.Context, request *web.Request) (*domain.Domain, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for AddProduct")
	}

	var r0 *domain.Domain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *web.Request) (*domain.Domain, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *web.Request) *domain.Domain); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Domain)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *web.Request) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteProduct provides a mock function with given fields: ctx, id
func (_m *Service) DeleteProduct(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteProduct")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetProducts provides a mock function with given fields: ctx
func (_m *Service) GetProducts(ctx context.Context) ([]*domain.Domain, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetProducts")
	}

	var r0 []*domain.Domain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*domain.Domain, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*domain.Domain); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Domain)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProduct provides a mock function with given fields: ctx, request, id
func (_m *Service) UpdateProduct(ctx context.Context, request *domain.Domain, id string) (*domain.Domain, error) {
	ret := _m.Called(ctx, request, id)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProduct")
	}

	var r0 *domain.Domain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Domain, string) (*domain.Domain, error)); ok {
		return rf(ctx, request, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Domain, string) *domain.Domain); ok {
		r0 = rf(ctx, request, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Domain)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.Domain, string) error); ok {
		r1 = rf(ctx, request, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
