// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ava-labs/avalanchego/vms/platformvm/state (interfaces: Chain)

// Package state is a generated GoMock package.
package state

import (
	big "math/big"
	reflect "reflect"
	time "time"

	ids "github.com/ava-labs/avalanchego/ids"
	avax "github.com/ava-labs/avalanchego/vms/components/avax"
	fx "github.com/ava-labs/avalanchego/vms/platformvm/fx"
	status "github.com/ava-labs/avalanchego/vms/platformvm/status"
	txs "github.com/ava-labs/avalanchego/vms/platformvm/txs"
	gomock "go.uber.org/mock/gomock"
)

// MockChain is a mock of Chain interface.
type MockChain struct {
	ctrl     *gomock.Controller
	recorder *MockChainMockRecorder
}

// MockChainMockRecorder is the mock recorder for MockChain.
type MockChainMockRecorder struct {
	mock *MockChain
}

// NewMockChain creates a new mock instance.
func NewMockChain(ctrl *gomock.Controller) *MockChain {
	mock := &MockChain{ctrl: ctrl}
	mock.recorder = &MockChainMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChain) EXPECT() *MockChainMockRecorder {
	return m.recorder
}

// AddChain mocks base method.
func (m *MockChain) AddChain(arg0 *txs.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddChain", arg0)
}

// AddChain indicates an expected call of AddChain.
func (mr *MockChainMockRecorder) AddChain(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddChain", reflect.TypeOf((*MockChain)(nil).AddChain), arg0)
}

// AddRewardUTXO mocks base method.
func (m *MockChain) AddRewardUTXO(arg0 ids.ID, arg1 *avax.UTXO) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddRewardUTXO", arg0, arg1)
}

// AddRewardUTXO indicates an expected call of AddRewardUTXO.
func (mr *MockChainMockRecorder) AddRewardUTXO(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRewardUTXO", reflect.TypeOf((*MockChain)(nil).AddRewardUTXO), arg0, arg1)
}

// AddSubnet mocks base method.
func (m *MockChain) AddSubnet(arg0 *txs.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddSubnet", arg0)
}

// AddSubnet indicates an expected call of AddSubnet.
func (mr *MockChainMockRecorder) AddSubnet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSubnet", reflect.TypeOf((*MockChain)(nil).AddSubnet), arg0)
}

// AddSubnetTransformation mocks base method.
func (m *MockChain) AddSubnetTransformation(arg0 *txs.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddSubnetTransformation", arg0)
}

// AddSubnetTransformation indicates an expected call of AddSubnetTransformation.
func (mr *MockChainMockRecorder) AddSubnetTransformation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSubnetTransformation", reflect.TypeOf((*MockChain)(nil).AddSubnetTransformation), arg0)
}

// AddTx mocks base method.
func (m *MockChain) AddTx(arg0 *txs.Tx, arg1 status.Status) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddTx", arg0, arg1)
}

// AddTx indicates an expected call of AddTx.
func (mr *MockChainMockRecorder) AddTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTx", reflect.TypeOf((*MockChain)(nil).AddTx), arg0, arg1)
}

// AddUTXO mocks base method.
func (m *MockChain) AddUTXO(arg0 *avax.UTXO) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddUTXO", arg0)
}

// AddUTXO indicates an expected call of AddUTXO.
func (mr *MockChainMockRecorder) AddUTXO(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUTXO", reflect.TypeOf((*MockChain)(nil).AddUTXO), arg0)
}

// DeleteCurrentDelegator mocks base method.
func (m *MockChain) DeleteCurrentDelegator(arg0 *Staker) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteCurrentDelegator", arg0)
}

// DeleteCurrentDelegator indicates an expected call of DeleteCurrentDelegator.
func (mr *MockChainMockRecorder) DeleteCurrentDelegator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCurrentDelegator", reflect.TypeOf((*MockChain)(nil).DeleteCurrentDelegator), arg0)
}

// DeleteCurrentValidator mocks base method.
func (m *MockChain) DeleteCurrentValidator(arg0 *Staker) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteCurrentValidator", arg0)
}

// DeleteCurrentValidator indicates an expected call of DeleteCurrentValidator.
func (mr *MockChainMockRecorder) DeleteCurrentValidator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCurrentValidator", reflect.TypeOf((*MockChain)(nil).DeleteCurrentValidator), arg0)
}

// DeletePendingDelegator mocks base method.
func (m *MockChain) DeletePendingDelegator(arg0 *Staker) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeletePendingDelegator", arg0)
}

// DeletePendingDelegator indicates an expected call of DeletePendingDelegator.
func (mr *MockChainMockRecorder) DeletePendingDelegator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePendingDelegator", reflect.TypeOf((*MockChain)(nil).DeletePendingDelegator), arg0)
}

// DeletePendingValidator mocks base method.
func (m *MockChain) DeletePendingValidator(arg0 *Staker) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeletePendingValidator", arg0)
}

// DeletePendingValidator indicates an expected call of DeletePendingValidator.
func (mr *MockChainMockRecorder) DeletePendingValidator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePendingValidator", reflect.TypeOf((*MockChain)(nil).DeletePendingValidator), arg0)
}

// DeleteUTXO mocks base method.
func (m *MockChain) DeleteUTXO(arg0 ids.ID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteUTXO", arg0)
}

// DeleteUTXO indicates an expected call of DeleteUTXO.
func (mr *MockChainMockRecorder) DeleteUTXO(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUTXO", reflect.TypeOf((*MockChain)(nil).DeleteUTXO), arg0)
}

// GetChains mocks base method.
func (m *MockChain) GetChains(arg0 ids.ID) ([]*txs.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChains", arg0)
	ret0, _ := ret[0].([]*txs.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChains indicates an expected call of GetChains.
func (mr *MockChainMockRecorder) GetChains(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChains", reflect.TypeOf((*MockChain)(nil).GetChains), arg0)
}

// GetCurrentDelegatorIterator mocks base method.
func (m *MockChain) GetCurrentDelegatorIterator(arg0 ids.ID, arg1 ids.NodeID) (StakerIterator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentDelegatorIterator", arg0, arg1)
	ret0, _ := ret[0].(StakerIterator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentDelegatorIterator indicates an expected call of GetCurrentDelegatorIterator.
func (mr *MockChainMockRecorder) GetCurrentDelegatorIterator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentDelegatorIterator", reflect.TypeOf((*MockChain)(nil).GetCurrentDelegatorIterator), arg0, arg1)
}

// GetCurrentStakerIterator mocks base method.
func (m *MockChain) GetCurrentStakerIterator() (StakerIterator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentStakerIterator")
	ret0, _ := ret[0].(StakerIterator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentStakerIterator indicates an expected call of GetCurrentStakerIterator.
func (mr *MockChainMockRecorder) GetCurrentStakerIterator() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentStakerIterator", reflect.TypeOf((*MockChain)(nil).GetCurrentStakerIterator))
}

// GetCurrentStakersLen mocks base method.
func (m *MockChain) GetCurrentStakersLen() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentStakersLen")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentStakersLen indicates an expected call of GetCurrentStakersLen.
func (mr *MockChainMockRecorder) GetCurrentStakersLen() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentStakersLen", reflect.TypeOf((*MockChain)(nil).GetCurrentStakersLen))
}

// GetCurrentSupply mocks base method.
func (m *MockChain) GetCurrentSupply(arg0 ids.ID) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentSupply", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentSupply indicates an expected call of GetCurrentSupply.
func (mr *MockChainMockRecorder) GetCurrentSupply(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentSupply", reflect.TypeOf((*MockChain)(nil).GetCurrentSupply), arg0)
}

// GetCurrentValidator mocks base method.
func (m *MockChain) GetCurrentValidator(arg0 ids.ID, arg1 ids.NodeID) (*Staker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentValidator", arg0, arg1)
	ret0, _ := ret[0].(*Staker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentValidator indicates an expected call of GetCurrentValidator.
func (mr *MockChainMockRecorder) GetCurrentValidator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentValidator", reflect.TypeOf((*MockChain)(nil).GetCurrentValidator), arg0, arg1)
}

// GetDelegateeReward mocks base method.
func (m *MockChain) GetDelegateeReward(arg0 ids.ID, arg1 ids.NodeID) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDelegateeReward", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDelegateeReward indicates an expected call of GetDelegateeReward.
func (mr *MockChainMockRecorder) GetDelegateeReward(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDelegateeReward", reflect.TypeOf((*MockChain)(nil).GetDelegateeReward), arg0, arg1)
}

// GetFeePerWeightStored mocks base method.
func (m *MockChain) GetFeePerWeightStored() (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeePerWeightStored")
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeePerWeightStored indicates an expected call of GetFeePerWeightStored.
func (mr *MockChainMockRecorder) GetFeePerWeightStored() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeePerWeightStored", reflect.TypeOf((*MockChain)(nil).GetFeePerWeightStored))
}

// GetLastAccumulatedFee mocks base method.
func (m *MockChain) GetLastAccumulatedFee() (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastAccumulatedFee")
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastAccumulatedFee indicates an expected call of GetLastAccumulatedFee.
func (mr *MockChainMockRecorder) GetLastAccumulatedFee() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastAccumulatedFee", reflect.TypeOf((*MockChain)(nil).GetLastAccumulatedFee))
}

// GetPendingDelegatorIterator mocks base method.
func (m *MockChain) GetPendingDelegatorIterator(arg0 ids.ID, arg1 ids.NodeID) (StakerIterator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPendingDelegatorIterator", arg0, arg1)
	ret0, _ := ret[0].(StakerIterator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPendingDelegatorIterator indicates an expected call of GetPendingDelegatorIterator.
func (mr *MockChainMockRecorder) GetPendingDelegatorIterator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPendingDelegatorIterator", reflect.TypeOf((*MockChain)(nil).GetPendingDelegatorIterator), arg0, arg1)
}

// GetPendingStakerIterator mocks base method.
func (m *MockChain) GetPendingStakerIterator() (StakerIterator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPendingStakerIterator")
	ret0, _ := ret[0].(StakerIterator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPendingStakerIterator indicates an expected call of GetPendingStakerIterator.
func (mr *MockChainMockRecorder) GetPendingStakerIterator() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPendingStakerIterator", reflect.TypeOf((*MockChain)(nil).GetPendingStakerIterator))
}

// GetPendingStakersLen mocks base method.
func (m *MockChain) GetPendingStakersLen() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPendingStakersLen")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPendingStakersLen indicates an expected call of GetPendingStakersLen.
func (mr *MockChainMockRecorder) GetPendingStakersLen() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPendingStakersLen", reflect.TypeOf((*MockChain)(nil).GetPendingStakersLen))
}

// GetPendingValidator mocks base method.
func (m *MockChain) GetPendingValidator(arg0 ids.ID, arg1 ids.NodeID) (*Staker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPendingValidator", arg0, arg1)
	ret0, _ := ret[0].(*Staker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPendingValidator indicates an expected call of GetPendingValidator.
func (mr *MockChainMockRecorder) GetPendingValidator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPendingValidator", reflect.TypeOf((*MockChain)(nil).GetPendingValidator), arg0, arg1)
}

// GetRewardUTXOs mocks base method.
func (m *MockChain) GetRewardUTXOs(arg0 ids.ID) ([]*avax.UTXO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRewardUTXOs", arg0)
	ret0, _ := ret[0].([]*avax.UTXO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRewardUTXOs indicates an expected call of GetRewardUTXOs.
func (mr *MockChainMockRecorder) GetRewardUTXOs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRewardUTXOs", reflect.TypeOf((*MockChain)(nil).GetRewardUTXOs), arg0)
}

// GetStakeSyncTimestamp mocks base method.
func (m *MockChain) GetStakeSyncTimestamp() (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStakeSyncTimestamp")
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStakeSyncTimestamp indicates an expected call of GetStakeSyncTimestamp.
func (mr *MockChainMockRecorder) GetStakeSyncTimestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStakeSyncTimestamp", reflect.TypeOf((*MockChain)(nil).GetStakeSyncTimestamp))
}

// GetStakerAccumulatedMintRate mocks base method.
func (m *MockChain) GetStakerAccumulatedMintRate() (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStakerAccumulatedMintRate")
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStakerAccumulatedMintRate indicates an expected call of GetStakerAccumulatedMintRate.
func (mr *MockChainMockRecorder) GetStakerAccumulatedMintRate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStakerAccumulatedMintRate", reflect.TypeOf((*MockChain)(nil).GetStakerAccumulatedMintRate))
}

// GetSubnetOwner mocks base method.
func (m *MockChain) GetSubnetOwner(arg0 ids.ID) (fx.Owner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubnetOwner", arg0)
	ret0, _ := ret[0].(fx.Owner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubnetOwner indicates an expected call of GetSubnetOwner.
func (mr *MockChainMockRecorder) GetSubnetOwner(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubnetOwner", reflect.TypeOf((*MockChain)(nil).GetSubnetOwner), arg0)
}

// GetSubnetTransformation mocks base method.
func (m *MockChain) GetSubnetTransformation(arg0 ids.ID) (*txs.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubnetTransformation", arg0)
	ret0, _ := ret[0].(*txs.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubnetTransformation indicates an expected call of GetSubnetTransformation.
func (mr *MockChainMockRecorder) GetSubnetTransformation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubnetTransformation", reflect.TypeOf((*MockChain)(nil).GetSubnetTransformation), arg0)
}

// GetSubnets mocks base method.
func (m *MockChain) GetSubnets() ([]*txs.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubnets")
	ret0, _ := ret[0].([]*txs.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubnets indicates an expected call of GetSubnets.
func (mr *MockChainMockRecorder) GetSubnets() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubnets", reflect.TypeOf((*MockChain)(nil).GetSubnets))
}

// GetTimestamp mocks base method.
func (m *MockChain) GetTimestamp() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTimestamp")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// GetTimestamp indicates an expected call of GetTimestamp.
func (mr *MockChainMockRecorder) GetTimestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTimestamp", reflect.TypeOf((*MockChain)(nil).GetTimestamp))
}

// GetTx mocks base method.
func (m *MockChain) GetTx(arg0 ids.ID) (*txs.Tx, status.Status, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTx", arg0)
	ret0, _ := ret[0].(*txs.Tx)
	ret1, _ := ret[1].(status.Status)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetTx indicates an expected call of GetTx.
func (mr *MockChainMockRecorder) GetTx(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTx", reflect.TypeOf((*MockChain)(nil).GetTx), arg0)
}

// GetUTXO mocks base method.
func (m *MockChain) GetUTXO(arg0 ids.ID) (*avax.UTXO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUTXO", arg0)
	ret0, _ := ret[0].(*avax.UTXO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUTXO indicates an expected call of GetUTXO.
func (mr *MockChainMockRecorder) GetUTXO(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUTXO", reflect.TypeOf((*MockChain)(nil).GetUTXO), arg0)
}

// PutCurrentDelegator mocks base method.
func (m *MockChain) PutCurrentDelegator(arg0 *Staker) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PutCurrentDelegator", arg0)
}

// PutCurrentDelegator indicates an expected call of PutCurrentDelegator.
func (mr *MockChainMockRecorder) PutCurrentDelegator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutCurrentDelegator", reflect.TypeOf((*MockChain)(nil).PutCurrentDelegator), arg0)
}

// PutCurrentValidator mocks base method.
func (m *MockChain) PutCurrentValidator(arg0 *Staker) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PutCurrentValidator", arg0)
}

// PutCurrentValidator indicates an expected call of PutCurrentValidator.
func (mr *MockChainMockRecorder) PutCurrentValidator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutCurrentValidator", reflect.TypeOf((*MockChain)(nil).PutCurrentValidator), arg0)
}

// PutPendingDelegator mocks base method.
func (m *MockChain) PutPendingDelegator(arg0 *Staker) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PutPendingDelegator", arg0)
}

// PutPendingDelegator indicates an expected call of PutPendingDelegator.
func (mr *MockChainMockRecorder) PutPendingDelegator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutPendingDelegator", reflect.TypeOf((*MockChain)(nil).PutPendingDelegator), arg0)
}

// PutPendingValidator mocks base method.
func (m *MockChain) PutPendingValidator(arg0 *Staker) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PutPendingValidator", arg0)
}

// PutPendingValidator indicates an expected call of PutPendingValidator.
func (mr *MockChainMockRecorder) PutPendingValidator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutPendingValidator", reflect.TypeOf((*MockChain)(nil).PutPendingValidator), arg0)
}

// SetCurrentSupply mocks base method.
func (m *MockChain) SetCurrentSupply(arg0 ids.ID, arg1 uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetCurrentSupply", arg0, arg1)
}

// SetCurrentSupply indicates an expected call of SetCurrentSupply.
func (mr *MockChainMockRecorder) SetCurrentSupply(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCurrentSupply", reflect.TypeOf((*MockChain)(nil).SetCurrentSupply), arg0, arg1)
}

// SetDelegateeReward mocks base method.
func (m *MockChain) SetDelegateeReward(arg0 ids.ID, arg1 ids.NodeID, arg2 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDelegateeReward", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDelegateeReward indicates an expected call of SetDelegateeReward.
func (mr *MockChainMockRecorder) SetDelegateeReward(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDelegateeReward", reflect.TypeOf((*MockChain)(nil).SetDelegateeReward), arg0, arg1, arg2)
}

// SetFeePerWeightStored mocks base method.
func (m *MockChain) SetFeePerWeightStored(arg0 *big.Int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetFeePerWeightStored", arg0)
}

// SetFeePerWeightStored indicates an expected call of SetFeePerWeightStored.
func (mr *MockChainMockRecorder) SetFeePerWeightStored(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFeePerWeightStored", reflect.TypeOf((*MockChain)(nil).SetFeePerWeightStored), arg0)
}

// SetLastAccumulatedFee mocks base method.
func (m *MockChain) SetLastAccumulatedFee(arg0 *big.Int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLastAccumulatedFee", arg0)
}

// SetLastAccumulatedFee indicates an expected call of SetLastAccumulatedFee.
func (mr *MockChainMockRecorder) SetLastAccumulatedFee(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLastAccumulatedFee", reflect.TypeOf((*MockChain)(nil).SetLastAccumulatedFee), arg0)
}

// SetStakeSyncTimestamp mocks base method.
func (m *MockChain) SetStakeSyncTimestamp(arg0 time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStakeSyncTimestamp", arg0)
}

// SetStakeSyncTimestamp indicates an expected call of SetStakeSyncTimestamp.
func (mr *MockChainMockRecorder) SetStakeSyncTimestamp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStakeSyncTimestamp", reflect.TypeOf((*MockChain)(nil).SetStakeSyncTimestamp), arg0)
}

// SetStakerAccumulatedMintRate mocks base method.
func (m *MockChain) SetStakerAccumulatedMintRate(arg0 *big.Int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStakerAccumulatedMintRate", arg0)
}

// SetStakerAccumulatedMintRate indicates an expected call of SetStakerAccumulatedMintRate.
func (mr *MockChainMockRecorder) SetStakerAccumulatedMintRate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStakerAccumulatedMintRate", reflect.TypeOf((*MockChain)(nil).SetStakerAccumulatedMintRate), arg0)
}

// SetTimestamp mocks base method.
func (m *MockChain) SetTimestamp(arg0 time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTimestamp", arg0)
}

// SetTimestamp indicates an expected call of SetTimestamp.
func (mr *MockChainMockRecorder) SetTimestamp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTimestamp", reflect.TypeOf((*MockChain)(nil).SetTimestamp), arg0)
}
