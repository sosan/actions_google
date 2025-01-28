// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	models "actions_google/pkg/domain/models"

	mock "github.com/stretchr/testify/mock"
)

// ActionsService is an autogenerated mock type for the ActionsService type
type ActionsService struct {
	mock.Mock
}

// GetAllContentFromGoogleSheets provides a mock function with given fields: newAction
func (_m *ActionsService) GetAllContentFromGoogleSheets(newAction *models.RequestGoogleAction) *[]byte {
	ret := _m.Called(newAction)

	if len(ret) == 0 {
		panic("no return value specified for GetAllContentFromGoogleSheets")
	}

	var r0 *[]byte
	if rf, ok := ret.Get(0).(func(*models.RequestGoogleAction) *[]byte); ok {
		r0 = rf(newAction)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]byte)
		}
	}

	return r0
}

// GetGoogleSheetByID provides a mock function with given fields: newAction
func (_m *ActionsService) GetGoogleSheetByID(newAction *models.RequestGoogleAction) *[]byte {
	ret := _m.Called(newAction)

	if len(ret) == 0 {
		panic("no return value specified for GetGoogleSheetByID")
	}

	var r0 *[]byte
	if rf, ok := ret.Get(0).(func(*models.RequestGoogleAction) *[]byte); ok {
		r0 = rf(newAction)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]byte)
		}
	}

	return r0
}

// GetNotion provides a mock function with given fields: newAction
func (_m *ActionsService) GetNotion(newAction *models.RequestGoogleAction) *[]byte {
	ret := _m.Called(newAction)

	if len(ret) == 0 {
		panic("no return value specified for GetNotion")
	}

	var r0 *[]byte
	if rf, ok := ret.Get(0).(func(*models.RequestGoogleAction) *[]byte); ok {
		r0 = rf(newAction)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]byte)
		}
	}

	return r0
}

// NewActionsService creates a new instance of ActionsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewActionsService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ActionsService {
	mock := &ActionsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
