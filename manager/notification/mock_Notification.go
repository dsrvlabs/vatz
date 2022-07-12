// Code generated by mockery v2.13.1. DO NOT EDIT.

package notification

import (
	v1 "github.com/dsrvlabs/vatz-proto/plugin/v1"
	mock "github.com/stretchr/testify/mock"
)

// MockNotification is an autogenerated mock type for the Notification type
type MockNotification struct {
	mock.Mock
}

// GetNotifyInfo provides a mock function with given fields: response, pluginName, methodName
func (_m *MockNotification) GetNotifyInfo(response *v1.ExecuteResponse, pluginName string, methodName string) NotifyInfo {
	ret := _m.Called(response, pluginName, methodName)

	var r0 NotifyInfo
	if rf, ok := ret.Get(0).(func(*v1.ExecuteResponse, string, string) NotifyInfo); ok {
		r0 = rf(response, pluginName, methodName)
	} else {
		r0 = ret.Get(0).(NotifyInfo)
	}

	return r0
}

// SendDiscord provides a mock function with given fields: msg, webhook
func (_m *MockNotification) SendDiscord(msg ReqMsg, webhook string) error {
	ret := _m.Called(msg, webhook)

	var r0 error
	if rf, ok := ret.Get(0).(func(ReqMsg, string) error); ok {
		r0 = rf(msg, webhook)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendNotification provides a mock function with given fields: request
func (_m *MockNotification) SendNotification(request ReqMsg) error {
	ret := _m.Called(request)

	var r0 error
	if rf, ok := ret.Get(0).(func(ReqMsg) error); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockNotification interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockNotification creates a new instance of MockNotification. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockNotification(t mockConstructorTestingTNewMockNotification) *MockNotification {
	mock := &MockNotification{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
