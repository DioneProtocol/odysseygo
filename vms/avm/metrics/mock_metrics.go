// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DioneProtocol/odysseygo/vms/avm/metrics (interfaces: Metrics)

// Package metrics is a generated GoMock package.
package metrics

import (
	http "net/http"
	reflect "reflect"

	block "github.com/DioneProtocol/odysseygo/vms/avm/block"
	txs "github.com/DioneProtocol/odysseygo/vms/avm/txs"
	rpc "github.com/gorilla/rpc/v2"
	gomock "go.uber.org/mock/gomock"
)

// MockMetrics is a mock of Metrics interface.
type MockMetrics struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsMockRecorder
}

// MockMetricsMockRecorder is the mock recorder for MockMetrics.
type MockMetricsMockRecorder struct {
	mock *MockMetrics
}

// NewMockMetrics creates a new mock instance.
func NewMockMetrics(ctrl *gomock.Controller) *MockMetrics {
	mock := &MockMetrics{ctrl: ctrl}
	mock.recorder = &MockMetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetrics) EXPECT() *MockMetricsMockRecorder {
	return m.recorder
}

// AfterRequest mocks base method.
func (m *MockMetrics) AfterRequest(arg0 *rpc.RequestInfo) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AfterRequest", arg0)
}

// AfterRequest indicates an expected call of AfterRequest.
func (mr *MockMetricsMockRecorder) AfterRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AfterRequest", reflect.TypeOf((*MockMetrics)(nil).AfterRequest), arg0)
}

// IncTxRefreshHits mocks base method.
func (m *MockMetrics) IncTxRefreshHits() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IncTxRefreshHits")
}

// IncTxRefreshHits indicates an expected call of IncTxRefreshHits.
func (mr *MockMetricsMockRecorder) IncTxRefreshHits() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncTxRefreshHits", reflect.TypeOf((*MockMetrics)(nil).IncTxRefreshHits))
}

// IncTxRefreshMisses mocks base method.
func (m *MockMetrics) IncTxRefreshMisses() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IncTxRefreshMisses")
}

// IncTxRefreshMisses indicates an expected call of IncTxRefreshMisses.
func (mr *MockMetricsMockRecorder) IncTxRefreshMisses() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncTxRefreshMisses", reflect.TypeOf((*MockMetrics)(nil).IncTxRefreshMisses))
}

// IncTxRefreshes mocks base method.
func (m *MockMetrics) IncTxRefreshes() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IncTxRefreshes")
}

// IncTxRefreshes indicates an expected call of IncTxRefreshes.
func (mr *MockMetricsMockRecorder) IncTxRefreshes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncTxRefreshes", reflect.TypeOf((*MockMetrics)(nil).IncTxRefreshes))
}

// InterceptRequest mocks base method.
func (m *MockMetrics) InterceptRequest(arg0 *rpc.RequestInfo) *http.Request {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InterceptRequest", arg0)
	ret0, _ := ret[0].(*http.Request)
	return ret0
}

// InterceptRequest indicates an expected call of InterceptRequest.
func (mr *MockMetricsMockRecorder) InterceptRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InterceptRequest", reflect.TypeOf((*MockMetrics)(nil).InterceptRequest), arg0)
}

// MarkBlockAccepted mocks base method.
func (m *MockMetrics) MarkBlockAccepted(arg0 block.Block) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkBlockAccepted", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkBlockAccepted indicates an expected call of MarkBlockAccepted.
func (mr *MockMetricsMockRecorder) MarkBlockAccepted(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkBlockAccepted", reflect.TypeOf((*MockMetrics)(nil).MarkBlockAccepted), arg0)
}

// MarkTxAccepted mocks base method.
func (m *MockMetrics) MarkTxAccepted(arg0 *txs.Tx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkTxAccepted", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkTxAccepted indicates an expected call of MarkTxAccepted.
func (mr *MockMetricsMockRecorder) MarkTxAccepted(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkTxAccepted", reflect.TypeOf((*MockMetrics)(nil).MarkTxAccepted), arg0)
}
