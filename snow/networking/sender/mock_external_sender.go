// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DioneProtocol/odysseygo/snow/networking/sender (interfaces: ExternalSender)

// Package sender is a generated GoMock package.
package sender

import (
	reflect "reflect"

	ids "github.com/DioneProtocol/odysseygo/ids"
	message "github.com/DioneProtocol/odysseygo/message"
	subnets "github.com/DioneProtocol/odysseygo/subnets"
	set "github.com/DioneProtocol/odysseygo/utils/set"
	gomock "github.com/golang/mock/gomock"
)

// MockExternalSender is a mock of ExternalSender interface.
type MockExternalSender struct {
	ctrl     *gomock.Controller
	recorder *MockExternalSenderMockRecorder
}

// MockExternalSenderMockRecorder is the mock recorder for MockExternalSender.
type MockExternalSenderMockRecorder struct {
	mock *MockExternalSender
}

// NewMockExternalSender creates a new mock instance.
func NewMockExternalSender(ctrl *gomock.Controller) *MockExternalSender {
	mock := &MockExternalSender{ctrl: ctrl}
	mock.recorder = &MockExternalSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExternalSender) EXPECT() *MockExternalSenderMockRecorder {
	return m.recorder
}

// Gossip mocks base method.
func (m *MockExternalSender) Gossip(arg0 message.OutboundMessage, arg1 ids.ID, arg2, arg3, arg4 int, arg5 subnets.Allower) set.Set[ids.NodeID] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Gossip", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(set.Set[ids.NodeID])
	return ret0
}

// Gossip indicates an expected call of Gossip.
func (mr *MockExternalSenderMockRecorder) Gossip(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Gossip", reflect.TypeOf((*MockExternalSender)(nil).Gossip), arg0, arg1, arg2, arg3, arg4, arg5)
}

// Send mocks base method.
func (m *MockExternalSender) Send(arg0 message.OutboundMessage, arg1 set.Set[ids.NodeID], arg2 ids.ID, arg3 subnets.Allower) set.Set[ids.NodeID] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(set.Set[ids.NodeID])
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockExternalSenderMockRecorder) Send(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockExternalSender)(nil).Send), arg0, arg1, arg2, arg3)
}
