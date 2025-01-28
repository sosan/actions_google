// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// JWTGenerator is an autogenerated mock type for the JWTGenerator type
type JWTGenerator struct {
	mock.Mock
}

// GenerateActionUserAssertionJWT provides a mock function with given fields: duration
func (_m *JWTGenerator) GenerateActionUserAssertionJWT(duration time.Duration) (string, error) {
	ret := _m.Called(duration)

	if len(ret) == 0 {
		panic("no return value specified for GenerateActionUserAssertionJWT")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(time.Duration) (string, error)); ok {
		return rf(duration)
	}
	if rf, ok := ret.Get(0).(func(time.Duration) string); ok {
		r0 = rf(duration)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(time.Duration) error); ok {
		r1 = rf(duration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateAppInstrospectJWT provides a mock function with given fields: duration
func (_m *JWTGenerator) GenerateAppInstrospectJWT(duration time.Duration) (string, error) {
	ret := _m.Called(duration)

	if len(ret) == 0 {
		panic("no return value specified for GenerateAppInstrospectJWT")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(time.Duration) (string, error)); ok {
		return rf(duration)
	}
	if rf, ok := ret.Get(0).(func(time.Duration) string); ok {
		r0 = rf(duration)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(time.Duration) error); ok {
		r1 = rf(duration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewJWTGenerator creates a new instance of JWTGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewJWTGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *JWTGenerator {
	mock := &JWTGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
