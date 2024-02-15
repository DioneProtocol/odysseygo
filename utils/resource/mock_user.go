// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DioneProtocol/odysseygo/utils/resource (interfaces: User)

// Package resource is a generated GoMock package.
package resource

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// AvailableDiskBytes mocks base method.
func (m *MockUser) AvailableDiskBytes() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AvailableDiskBytes")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// AvailableDiskBytes indicates an expected call of AvailableDiskBytes.
func (mr *MockUserMockRecorder) AvailableDiskBytes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AvailableDiskBytes", reflect.TypeOf((*MockUser)(nil).AvailableDiskBytes))
}

// CPUUsage mocks base method.
func (m *MockUser) CPUUsage() float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CPUUsage")
	ret0, _ := ret[0].(float64)
	return ret0
}

// CPUUsage indicates an expected call of CPUUsage.
func (mr *MockUserMockRecorder) CPUUsage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CPUUsage", reflect.TypeOf((*MockUser)(nil).CPUUsage))
}

// DiskUsage mocks base method.
func (m *MockUser) DiskUsage() (float64, float64) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DiskUsage")
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(float64)
	return ret0, ret1
}

// DiskUsage indicates an expected call of DiskUsage.
func (mr *MockUserMockRecorder) DiskUsage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiskUsage", reflect.TypeOf((*MockUser)(nil).DiskUsage))
}
