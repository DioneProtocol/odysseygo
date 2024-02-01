// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ava-labs/avalanchego/vms/components/feecollector (interfaces: FeeCollector)

// Package feecollector is a generated GoMock package.
package feecollector

import (
	big "math/big"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockFeeCollector is a mock of FeeCollector interface.
type MockFeeCollector struct {
	ctrl     *gomock.Controller
	recorder *MockFeeCollectorMockRecorder
}

// MockFeeCollectorMockRecorder is the mock recorder for MockFeeCollector.
type MockFeeCollectorMockRecorder struct {
	mock *MockFeeCollector
}

// NewMockFeeCollector creates a new mock instance.
func NewMockFeeCollector(ctrl *gomock.Controller) *MockFeeCollector {
	mock := &MockFeeCollector{ctrl: ctrl}
	mock.recorder = &MockFeeCollectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFeeCollector) EXPECT() *MockFeeCollectorMockRecorder {
	return m.recorder
}

// AddCChainValue mocks base method.
func (m *MockFeeCollector) AddCChainValue(arg0 *big.Int) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCChainValue", arg0)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCChainValue indicates an expected call of AddCChainValue.
func (mr *MockFeeCollectorMockRecorder) AddCChainValue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCChainValue", reflect.TypeOf((*MockFeeCollector)(nil).AddCChainValue), arg0)
}

// AddPChainValue mocks base method.
func (m *MockFeeCollector) AddPChainValue(arg0 *big.Int) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPChainValue", arg0)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddPChainValue indicates an expected call of AddPChainValue.
func (mr *MockFeeCollectorMockRecorder) AddPChainValue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPChainValue", reflect.TypeOf((*MockFeeCollector)(nil).AddPChainValue), arg0)
}

// AddXChainValue mocks base method.
func (m *MockFeeCollector) AddXChainValue(arg0 *big.Int) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddXChainValue", arg0)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddXChainValue indicates an expected call of AddXChainValue.
func (mr *MockFeeCollectorMockRecorder) AddXChainValue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddXChainValue", reflect.TypeOf((*MockFeeCollector)(nil).AddXChainValue), arg0)
}

// GetCChainValue mocks base method.
func (m *MockFeeCollector) GetCChainValue() *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCChainValue")
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// GetCChainValue indicates an expected call of GetCChainValue.
func (mr *MockFeeCollectorMockRecorder) GetCChainValue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCChainValue", reflect.TypeOf((*MockFeeCollector)(nil).GetCChainValue))
}

// GetPChainValue mocks base method.
func (m *MockFeeCollector) GetPChainValue() *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPChainValue")
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// GetPChainValue indicates an expected call of GetPChainValue.
func (mr *MockFeeCollectorMockRecorder) GetPChainValue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPChainValue", reflect.TypeOf((*MockFeeCollector)(nil).GetPChainValue))
}

// GetXChainValue mocks base method.
func (m *MockFeeCollector) GetXChainValue() *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetXChainValue")
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// GetXChainValue indicates an expected call of GetXChainValue.
func (mr *MockFeeCollectorMockRecorder) GetXChainValue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetXChainValue", reflect.TypeOf((*MockFeeCollector)(nil).GetXChainValue))
}

// SubCChainValue mocks base method.
func (m *MockFeeCollector) SubCChainValue(arg0 *big.Int) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubCChainValue", arg0)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubCChainValue indicates an expected call of SubCChainValue.
func (mr *MockFeeCollectorMockRecorder) SubCChainValue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubCChainValue", reflect.TypeOf((*MockFeeCollector)(nil).SubCChainValue), arg0)
}

// SubPChainValue mocks base method.
func (m *MockFeeCollector) SubPChainValue(arg0 *big.Int) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubPChainValue", arg0)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubPChainValue indicates an expected call of SubPChainValue.
func (mr *MockFeeCollectorMockRecorder) SubPChainValue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubPChainValue", reflect.TypeOf((*MockFeeCollector)(nil).SubPChainValue), arg0)
}

// SubXChainValue mocks base method.
func (m *MockFeeCollector) SubXChainValue(arg0 *big.Int) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubXChainValue", arg0)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubXChainValue indicates an expected call of SubXChainValue.
func (mr *MockFeeCollectorMockRecorder) SubXChainValue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubXChainValue", reflect.TypeOf((*MockFeeCollector)(nil).SubXChainValue), arg0)
}