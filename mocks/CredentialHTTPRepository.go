// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	models "actions_google/pkg/domain/models"

	mock "github.com/stretchr/testify/mock"
)

// CredentialHTTPRepository is an autogenerated mock type for the CredentialHTTPRepository type
type CredentialHTTPRepository struct {
	mock.Mock
}

// GetAllCredentials provides a mock function with given fields: userID, limitCount
func (_m *CredentialHTTPRepository) GetAllCredentials(userID *string, limitCount uint64) (*[]models.RequestExchangeCredential, error) {
	ret := _m.Called(userID, limitCount)

	if len(ret) == 0 {
		panic("no return value specified for GetAllCredentials")
	}

	var r0 *[]models.RequestExchangeCredential
	var r1 error
	if rf, ok := ret.Get(0).(func(*string, uint64) (*[]models.RequestExchangeCredential, error)); ok {
		return rf(userID, limitCount)
	}
	if rf, ok := ret.Get(0).(func(*string, uint64) *[]models.RequestExchangeCredential); ok {
		r0 = rf(userID, limitCount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.RequestExchangeCredential)
		}
	}

	if rf, ok := ret.Get(1).(func(*string, uint64) error); ok {
		r1 = rf(userID, limitCount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCredentialByID provides a mock function with given fields: userID, credentialID, limitCount
func (_m *CredentialHTTPRepository) GetCredentialByID(userID *string, credentialID *string, limitCount uint64) (*models.RequestExchangeCredential, error) {
	ret := _m.Called(userID, credentialID, limitCount)

	if len(ret) == 0 {
		panic("no return value specified for GetCredentialByID")
	}

	var r0 *models.RequestExchangeCredential
	var r1 error
	if rf, ok := ret.Get(0).(func(*string, *string, uint64) (*models.RequestExchangeCredential, error)); ok {
		return rf(userID, credentialID, limitCount)
	}
	if rf, ok := ret.Get(0).(func(*string, *string, uint64) *models.RequestExchangeCredential); ok {
		r0 = rf(userID, credentialID, limitCount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.RequestExchangeCredential)
		}
	}

	if rf, ok := ret.Get(1).(func(*string, *string, uint64) error); ok {
		r1 = rf(userID, credentialID, limitCount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCredentialHTTPRepository creates a new instance of CredentialHTTPRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCredentialHTTPRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *CredentialHTTPRepository {
	mock := &CredentialHTTPRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
