// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ava-labs/avalanchego/x/merkledb (interfaces: MerkleDB)

// Package merkledb is a generated GoMock package.
package merkledb

import (
	context "context"
	reflect "reflect"

	database "github.com/ava-labs/avalanchego/database"
	ids "github.com/ava-labs/avalanchego/ids"
	maybe "github.com/ava-labs/avalanchego/utils/maybe"
	gomock "go.uber.org/mock/gomock"
)

// MockMerkleDB is a mock of MerkleDB interface.
type MockMerkleDB struct {
	ctrl     *gomock.Controller
	recorder *MockMerkleDBMockRecorder
}

// MockMerkleDBMockRecorder is the mock recorder for MockMerkleDB.
type MockMerkleDBMockRecorder struct {
	mock *MockMerkleDB
}

// NewMockMerkleDB creates a new mock instance.
func NewMockMerkleDB(ctrl *gomock.Controller) *MockMerkleDB {
	mock := &MockMerkleDB{ctrl: ctrl}
	mock.recorder = &MockMerkleDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMerkleDB) EXPECT() *MockMerkleDBMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockMerkleDB) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockMerkleDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockMerkleDB)(nil).Close))
}

// CommitChangeProof mocks base method.
func (m *MockMerkleDB) CommitChangeProof(arg0 context.Context, arg1 *ChangeProof) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitChangeProof", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitChangeProof indicates an expected call of CommitChangeProof.
func (mr *MockMerkleDBMockRecorder) CommitChangeProof(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitChangeProof", reflect.TypeOf((*MockMerkleDB)(nil).CommitChangeProof), arg0, arg1)
}

// CommitRangeProof mocks base method.
func (m *MockMerkleDB) CommitRangeProof(arg0 context.Context, arg1 maybe.Maybe[[]uint8], arg2 *RangeProof) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitRangeProof", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitRangeProof indicates an expected call of CommitRangeProof.
func (mr *MockMerkleDBMockRecorder) CommitRangeProof(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitRangeProof", reflect.TypeOf((*MockMerkleDB)(nil).CommitRangeProof), arg0, arg1, arg2)
}

// Compact mocks base method.
func (m *MockMerkleDB) Compact(arg0, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compact", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Compact indicates an expected call of Compact.
func (mr *MockMerkleDBMockRecorder) Compact(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compact", reflect.TypeOf((*MockMerkleDB)(nil).Compact), arg0, arg1)
}

// Delete mocks base method.
func (m *MockMerkleDB) Delete(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockMerkleDBMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockMerkleDB)(nil).Delete), arg0)
}

// Get mocks base method.
func (m *MockMerkleDB) Get(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockMerkleDBMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMerkleDB)(nil).Get), arg0)
}

// GetChangeProof mocks base method.
func (m *MockMerkleDB) GetChangeProof(arg0 context.Context, arg1, arg2 ids.ID, arg3, arg4 maybe.Maybe[[]uint8], arg5 int) (*ChangeProof, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChangeProof", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*ChangeProof)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChangeProof indicates an expected call of GetChangeProof.
func (mr *MockMerkleDBMockRecorder) GetChangeProof(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChangeProof", reflect.TypeOf((*MockMerkleDB)(nil).GetChangeProof), arg0, arg1, arg2, arg3, arg4, arg5)
}

// GetMerkleRoot mocks base method.
func (m *MockMerkleDB) GetMerkleRoot(arg0 context.Context) (ids.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMerkleRoot", arg0)
	ret0, _ := ret[0].(ids.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMerkleRoot indicates an expected call of GetMerkleRoot.
func (mr *MockMerkleDBMockRecorder) GetMerkleRoot(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMerkleRoot", reflect.TypeOf((*MockMerkleDB)(nil).GetMerkleRoot), arg0)
}

// GetProof mocks base method.
func (m *MockMerkleDB) GetProof(arg0 context.Context, arg1 []byte) (*Proof, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProof", arg0, arg1)
	ret0, _ := ret[0].(*Proof)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProof indicates an expected call of GetProof.
func (mr *MockMerkleDBMockRecorder) GetProof(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProof", reflect.TypeOf((*MockMerkleDB)(nil).GetProof), arg0, arg1)
}

// GetRangeProof mocks base method.
func (m *MockMerkleDB) GetRangeProof(arg0 context.Context, arg1, arg2 maybe.Maybe[[]uint8], arg3 int) (*RangeProof, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRangeProof", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*RangeProof)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRangeProof indicates an expected call of GetRangeProof.
func (mr *MockMerkleDBMockRecorder) GetRangeProof(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRangeProof", reflect.TypeOf((*MockMerkleDB)(nil).GetRangeProof), arg0, arg1, arg2, arg3)
}

// GetRangeProofAtRoot mocks base method.
func (m *MockMerkleDB) GetRangeProofAtRoot(arg0 context.Context, arg1 ids.ID, arg2, arg3 maybe.Maybe[[]uint8], arg4 int) (*RangeProof, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRangeProofAtRoot", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*RangeProof)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRangeProofAtRoot indicates an expected call of GetRangeProofAtRoot.
func (mr *MockMerkleDBMockRecorder) GetRangeProofAtRoot(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRangeProofAtRoot", reflect.TypeOf((*MockMerkleDB)(nil).GetRangeProofAtRoot), arg0, arg1, arg2, arg3, arg4)
}

// GetValue mocks base method.
func (m *MockMerkleDB) GetValue(arg0 context.Context, arg1 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValue", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValue indicates an expected call of GetValue.
func (mr *MockMerkleDBMockRecorder) GetValue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValue", reflect.TypeOf((*MockMerkleDB)(nil).GetValue), arg0, arg1)
}

// GetValues mocks base method.
func (m *MockMerkleDB) GetValues(arg0 context.Context, arg1 [][]byte) ([][]byte, []error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValues", arg0, arg1)
	ret0, _ := ret[0].([][]byte)
	ret1, _ := ret[1].([]error)
	return ret0, ret1
}

// GetValues indicates an expected call of GetValues.
func (mr *MockMerkleDBMockRecorder) GetValues(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValues", reflect.TypeOf((*MockMerkleDB)(nil).GetValues), arg0, arg1)
}

// Has mocks base method.
func (m *MockMerkleDB) Has(arg0 []byte) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Has", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Has indicates an expected call of Has.
func (mr *MockMerkleDBMockRecorder) Has(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Has", reflect.TypeOf((*MockMerkleDB)(nil).Has), arg0)
}

// HealthCheck mocks base method.
func (m *MockMerkleDB) HealthCheck(arg0 context.Context) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthCheck", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HealthCheck indicates an expected call of HealthCheck.
func (mr *MockMerkleDBMockRecorder) HealthCheck(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthCheck", reflect.TypeOf((*MockMerkleDB)(nil).HealthCheck), arg0)
}

// NewBatch mocks base method.
func (m *MockMerkleDB) NewBatch() database.Batch {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewBatch")
	ret0, _ := ret[0].(database.Batch)
	return ret0
}

// NewBatch indicates an expected call of NewBatch.
func (mr *MockMerkleDBMockRecorder) NewBatch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewBatch", reflect.TypeOf((*MockMerkleDB)(nil).NewBatch))
}

// NewIterator mocks base method.
func (m *MockMerkleDB) NewIterator() database.Iterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewIterator")
	ret0, _ := ret[0].(database.Iterator)
	return ret0
}

// NewIterator indicates an expected call of NewIterator.
func (mr *MockMerkleDBMockRecorder) NewIterator() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewIterator", reflect.TypeOf((*MockMerkleDB)(nil).NewIterator))
}

// NewIteratorWithPrefix mocks base method.
func (m *MockMerkleDB) NewIteratorWithPrefix(arg0 []byte) database.Iterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewIteratorWithPrefix", arg0)
	ret0, _ := ret[0].(database.Iterator)
	return ret0
}

// NewIteratorWithPrefix indicates an expected call of NewIteratorWithPrefix.
func (mr *MockMerkleDBMockRecorder) NewIteratorWithPrefix(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewIteratorWithPrefix", reflect.TypeOf((*MockMerkleDB)(nil).NewIteratorWithPrefix), arg0)
}

// NewIteratorWithStart mocks base method.
func (m *MockMerkleDB) NewIteratorWithStart(arg0 []byte) database.Iterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewIteratorWithStart", arg0)
	ret0, _ := ret[0].(database.Iterator)
	return ret0
}

// NewIteratorWithStart indicates an expected call of NewIteratorWithStart.
func (mr *MockMerkleDBMockRecorder) NewIteratorWithStart(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewIteratorWithStart", reflect.TypeOf((*MockMerkleDB)(nil).NewIteratorWithStart), arg0)
}

// NewIteratorWithStartAndPrefix mocks base method.
func (m *MockMerkleDB) NewIteratorWithStartAndPrefix(arg0, arg1 []byte) database.Iterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewIteratorWithStartAndPrefix", arg0, arg1)
	ret0, _ := ret[0].(database.Iterator)
	return ret0
}

// NewIteratorWithStartAndPrefix indicates an expected call of NewIteratorWithStartAndPrefix.
func (mr *MockMerkleDBMockRecorder) NewIteratorWithStartAndPrefix(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewIteratorWithStartAndPrefix", reflect.TypeOf((*MockMerkleDB)(nil).NewIteratorWithStartAndPrefix), arg0, arg1)
}

// NewView mocks base method.
func (m *MockMerkleDB) NewView(arg0 []database.BatchOp) (TrieView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewView", arg0)
	ret0, _ := ret[0].(TrieView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewView indicates an expected call of NewView.
func (mr *MockMerkleDBMockRecorder) NewView(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewView", reflect.TypeOf((*MockMerkleDB)(nil).NewView), arg0)
}

// Put mocks base method.
func (m *MockMerkleDB) Put(arg0, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockMerkleDBMockRecorder) Put(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockMerkleDB)(nil).Put), arg0, arg1)
}

// VerifyChangeProof mocks base method.
func (m *MockMerkleDB) VerifyChangeProof(arg0 context.Context, arg1 *ChangeProof, arg2, arg3 maybe.Maybe[[]uint8], arg4 ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyChangeProof", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyChangeProof indicates an expected call of VerifyChangeProof.
func (mr *MockMerkleDBMockRecorder) VerifyChangeProof(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyChangeProof", reflect.TypeOf((*MockMerkleDB)(nil).VerifyChangeProof), arg0, arg1, arg2, arg3, arg4)
}

// getEditableNode mocks base method.
func (m *MockMerkleDB) getEditableNode(arg0 path) (*node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getEditableNode", arg0)
	ret0, _ := ret[0].(*node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getEditableNode indicates an expected call of getEditableNode.
func (mr *MockMerkleDBMockRecorder) getEditableNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getEditableNode", reflect.TypeOf((*MockMerkleDB)(nil).getEditableNode), arg0)
}

// getValue mocks base method.
func (m *MockMerkleDB) getValue(arg0 path, arg1 bool) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getValue", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getValue indicates an expected call of getValue.
func (mr *MockMerkleDBMockRecorder) getValue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getValue", reflect.TypeOf((*MockMerkleDB)(nil).getValue), arg0, arg1)
}
