// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DioneProtocol/odysseygo/snow/validators (interfaces: State)

// Package validators is a generated GoMock package.
package validators

import (
	context "context"
	reflect "reflect"

	ids "github.com/DioneProtocol/odysseygo/ids"
	gomock "github.com/golang/mock/gomock"
)

// MockState is a mock of State interface.
type MockState struct {
	ctrl     *gomock.Controller
	recorder *MockStateMockRecorder
}

// MockStateMockRecorder is the mock recorder for MockState.
type MockStateMockRecorder struct {
	mock *MockState
}

// NewMockState creates a new mock instance.
func NewMockState(ctrl *gomock.Controller) *MockState {
	mock := &MockState{ctrl: ctrl}
	mock.recorder = &MockStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockState) EXPECT() *MockStateMockRecorder {
	return m.recorder
}

// GetCurrentHeight mocks base method.
func (m *MockState) GetCurrentHeight(arg0 context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentHeight", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentHeight indicates an expected call of GetCurrentHeight.
func (mr *MockStateMockRecorder) GetCurrentHeight(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentHeight", reflect.TypeOf((*MockState)(nil).GetCurrentHeight), arg0)
}

// GetMinimumHeight mocks base method.
func (m *MockState) GetMinimumHeight(arg0 context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMinimumHeight", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMinimumHeight indicates an expected call of GetMinimumHeight.
func (mr *MockStateMockRecorder) GetMinimumHeight(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMinimumHeight", reflect.TypeOf((*MockState)(nil).GetMinimumHeight), arg0)
}

// GetSubnetID mocks base method.
func (m *MockState) GetSubnetID(arg0 context.Context, arg1 ids.ID) (ids.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubnetID", arg0, arg1)
	ret0, _ := ret[0].(ids.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubnetID indicates an expected call of GetSubnetID.
func (mr *MockStateMockRecorder) GetSubnetID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubnetID", reflect.TypeOf((*MockState)(nil).GetSubnetID), arg0, arg1)
}

// GetValidatorSet mocks base method.
func (m *MockState) GetValidatorSet(arg0 context.Context, arg1 uint64, arg2 ids.ID) (map[ids.NodeID]*GetValidatorOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorSet", arg0, arg1, arg2)
	ret0, _ := ret[0].(map[ids.NodeID]*GetValidatorOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidatorSet indicates an expected call of GetValidatorSet.
func (mr *MockStateMockRecorder) GetValidatorSet(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorSet", reflect.TypeOf((*MockState)(nil).GetValidatorSet), arg0, arg1, arg2)
}
