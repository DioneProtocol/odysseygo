// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DioneProtocol/odysseygo/snow/networking/timeout (interfaces: Manager)

// Package timeout is a generated GoMock package.
package timeout

import (
	reflect "reflect"
	time "time"

	ids "github.com/DioneProtocol/odysseygo/ids"
	message "github.com/DioneProtocol/odysseygo/message"
	snow "github.com/DioneProtocol/odysseygo/snow"
	gomock "go.uber.org/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Dispatch mocks base method.
func (m *MockManager) Dispatch() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dispatch")
}

// Dispatch indicates an expected call of Dispatch.
func (mr *MockManagerMockRecorder) Dispatch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispatch", reflect.TypeOf((*MockManager)(nil).Dispatch))
}

// IsBenched mocks base method.
func (m *MockManager) IsBenched(arg0 ids.NodeID, arg1 ids.ID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsBenched", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsBenched indicates an expected call of IsBenched.
func (mr *MockManagerMockRecorder) IsBenched(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsBenched", reflect.TypeOf((*MockManager)(nil).IsBenched), arg0, arg1)
}

// RegisterChain mocks base method.
func (m *MockManager) RegisterChain(arg0 *snow.ConsensusContext) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterChain", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterChain indicates an expected call of RegisterChain.
func (mr *MockManagerMockRecorder) RegisterChain(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterChain", reflect.TypeOf((*MockManager)(nil).RegisterChain), arg0)
}

// RegisterRequest mocks base method.
func (m *MockManager) RegisterRequest(arg0 ids.NodeID, arg1 ids.ID, arg2 bool, arg3 ids.RequestID, arg4 func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterRequest", arg0, arg1, arg2, arg3, arg4)
}

// RegisterRequest indicates an expected call of RegisterRequest.
func (mr *MockManagerMockRecorder) RegisterRequest(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterRequest", reflect.TypeOf((*MockManager)(nil).RegisterRequest), arg0, arg1, arg2, arg3, arg4)
}

// RegisterRequestToUnreachableValidator mocks base method.
func (m *MockManager) RegisterRequestToUnreachableValidator() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterRequestToUnreachableValidator")
}

// RegisterRequestToUnreachableValidator indicates an expected call of RegisterRequestToUnreachableValidator.
func (mr *MockManagerMockRecorder) RegisterRequestToUnreachableValidator() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterRequestToUnreachableValidator", reflect.TypeOf((*MockManager)(nil).RegisterRequestToUnreachableValidator))
}

// RegisterResponse mocks base method.
func (m *MockManager) RegisterResponse(arg0 ids.NodeID, arg1 ids.ID, arg2 ids.RequestID, arg3 message.Op, arg4 time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterResponse", arg0, arg1, arg2, arg3, arg4)
}

// RegisterResponse indicates an expected call of RegisterResponse.
func (mr *MockManagerMockRecorder) RegisterResponse(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterResponse", reflect.TypeOf((*MockManager)(nil).RegisterResponse), arg0, arg1, arg2, arg3, arg4)
}

// RemoveRequest mocks base method.
func (m *MockManager) RemoveRequest(arg0 ids.RequestID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveRequest", arg0)
}

// RemoveRequest indicates an expected call of RemoveRequest.
func (mr *MockManagerMockRecorder) RemoveRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRequest", reflect.TypeOf((*MockManager)(nil).RemoveRequest), arg0)
}

// TimeoutDuration mocks base method.
func (m *MockManager) TimeoutDuration() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TimeoutDuration")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// TimeoutDuration indicates an expected call of TimeoutDuration.
func (mr *MockManagerMockRecorder) TimeoutDuration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TimeoutDuration", reflect.TypeOf((*MockManager)(nil).TimeoutDuration))
}
