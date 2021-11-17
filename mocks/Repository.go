// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// ChangeBalance provides a mock function with given fields: userId, amount
func (_m *Repository) ChangeBalance(userId int64, amount float32) error {
	ret := _m.Called(userId, amount)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, float32) error); ok {
		r0 = rf(userId, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetBalance provides a mock function with given fields: userId
func (_m *Repository) GetBalance(userId int64) (float32, error) {
	ret := _m.Called(userId)

	var r0 float32
	if rf, ok := ret.Get(0).(func(int64) float32); ok {
		r0 = rf(userId)
	} else {
		r0 = ret.Get(0).(float32)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransferMoney provides a mock function with given fields: srcUserId, dstUserId, amount
func (_m *Repository) TransferMoney(srcUserId int64, dstUserId int64, amount float32) error {
	ret := _m.Called(srcUserId, dstUserId, amount)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, int64, float32) error); ok {
		r0 = rf(srcUserId, dstUserId, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}