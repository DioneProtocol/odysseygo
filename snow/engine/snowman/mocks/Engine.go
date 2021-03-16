// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	common "github.com/ava-labs/avalanchego/snow/engine/common"
	block "github.com/ava-labs/avalanchego/snow/engine/snowman/block"

	consensussnowman "github.com/ava-labs/avalanchego/snow/consensus/snowman"

	ids "github.com/ava-labs/avalanchego/ids"

	mock "github.com/stretchr/testify/mock"

	snow "github.com/ava-labs/avalanchego/snow"

	snowman "github.com/ava-labs/avalanchego/snow/engine/snowman"
)

// SnowmanEngine is an autogenerated mock type for the SnowmanEngine type
type SnowmanEngine struct {
	mock.Mock
}

// Accepted provides a mock function with given fields: validatorID, requestID, containerIDs
func (_m *SnowmanEngine) Accepted(validatorID ids.ShortID, requestID uint32, containerIDs []ids.ID) error {
	ret := _m.Called(validatorID, requestID, containerIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32, []ids.ID) error); ok {
		r0 = rf(validatorID, requestID, containerIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AcceptedFrontier provides a mock function with given fields: validatorID, requestID, containerIDs
func (_m *SnowmanEngine) AcceptedFrontier(validatorID ids.ShortID, requestID uint32, containerIDs []ids.ID) error {
	ret := _m.Called(validatorID, requestID, containerIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32, []ids.ID) error); ok {
		r0 = rf(validatorID, requestID, containerIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Chits provides a mock function with given fields: validatorID, requestID, containerIDs
func (_m *SnowmanEngine) Chits(validatorID ids.ShortID, requestID uint32, containerIDs []ids.ID) error {
	ret := _m.Called(validatorID, requestID, containerIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32, []ids.ID) error); ok {
		r0 = rf(validatorID, requestID, containerIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Connected provides a mock function with given fields: validatorID
func (_m *SnowmanEngine) Connected(validatorID ids.ShortID) error {
	ret := _m.Called(validatorID)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID) error); ok {
		r0 = rf(validatorID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Context provides a mock function with given fields:
func (_m *SnowmanEngine) Context() *snow.Context {
	ret := _m.Called()

	var r0 *snow.Context
	if rf, ok := ret.Get(0).(func() *snow.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*snow.Context)
		}
	}

	return r0
}

// Disconnected provides a mock function with given fields: validatorID
func (_m *SnowmanEngine) Disconnected(validatorID ids.ShortID) error {
	ret := _m.Called(validatorID)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID) error); ok {
		r0 = rf(validatorID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: validatorID, requestID, containerID
func (_m *SnowmanEngine) Get(validatorID ids.ShortID, requestID uint32, containerID ids.ID) error {
	ret := _m.Called(validatorID, requestID, containerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32, ids.ID) error); ok {
		r0 = rf(validatorID, requestID, containerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAccepted provides a mock function with given fields: validatorID, requestID, containerIDs
func (_m *SnowmanEngine) GetAccepted(validatorID ids.ShortID, requestID uint32, containerIDs []ids.ID) error {
	ret := _m.Called(validatorID, requestID, containerIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32, []ids.ID) error); ok {
		r0 = rf(validatorID, requestID, containerIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAcceptedFailed provides a mock function with given fields: validatorID, requestID
func (_m *SnowmanEngine) GetAcceptedFailed(validatorID ids.ShortID, requestID uint32) error {
	ret := _m.Called(validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32) error); ok {
		r0 = rf(validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAcceptedFrontier provides a mock function with given fields: validatorID, requestID
func (_m *SnowmanEngine) GetAcceptedFrontier(validatorID ids.ShortID, requestID uint32) error {
	ret := _m.Called(validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32) error); ok {
		r0 = rf(validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAcceptedFrontierFailed provides a mock function with given fields: validatorID, requestID
func (_m *SnowmanEngine) GetAcceptedFrontierFailed(validatorID ids.ShortID, requestID uint32) error {
	ret := _m.Called(validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32) error); ok {
		r0 = rf(validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAncestors provides a mock function with given fields: validatorID, requestID, containerID
func (_m *SnowmanEngine) GetAncestors(validatorID ids.ShortID, requestID uint32, containerID ids.ID) error {
	ret := _m.Called(validatorID, requestID, containerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32, ids.ID) error); ok {
		r0 = rf(validatorID, requestID, containerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAncestorsFailed provides a mock function with given fields: validatorID, requestID
func (_m *SnowmanEngine) GetAncestorsFailed(validatorID ids.ShortID, requestID uint32) error {
	ret := _m.Called(validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32) error); ok {
		r0 = rf(validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetBlock provides a mock function with given fields: blkID
func (_m *SnowmanEngine) GetBlock(blkID ids.ID) (consensussnowman.Block, error) {
	ret := _m.Called(blkID)

	var r0 consensussnowman.Block
	if rf, ok := ret.Get(0).(func(ids.ID) consensussnowman.Block); ok {
		r0 = rf(blkID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(consensussnowman.Block)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ids.ID) error); ok {
		r1 = rf(blkID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFailed provides a mock function with given fields: validatorID, requestID
func (_m *SnowmanEngine) GetFailed(validatorID ids.ShortID, requestID uint32) error {
	ret := _m.Called(validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32) error); ok {
		r0 = rf(validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetVM provides a mock function with given fields:
func (_m *SnowmanEngine) GetVM() block.ChainVM {
	ret := _m.Called()

	var r0 block.ChainVM
	if rf, ok := ret.Get(0).(func() block.ChainVM); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(block.ChainVM)
		}
	}

	return r0
}

// Gossip provides a mock function with given fields:
func (_m *SnowmanEngine) Gossip() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HealthCheck provides a mock function with given fields:
func (_m *SnowmanEngine) HealthCheck() (interface{}, error) {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Initialize provides a mock function with given fields: _a0
func (_m *SnowmanEngine) Initialize(_a0 snowman.Config) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(snowman.Config) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsBootstrapped provides a mock function with given fields:
func (_m *SnowmanEngine) IsBootstrapped() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MultiPut provides a mock function with given fields: validatorID, requestID, containers
func (_m *SnowmanEngine) MultiPut(validatorID ids.ShortID, requestID uint32, containers [][]byte) error {
	ret := _m.Called(validatorID, requestID, containers)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32, [][]byte) error); ok {
		r0 = rf(validatorID, requestID, containers)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Notify provides a mock function with given fields: _a0
func (_m *SnowmanEngine) Notify(_a0 common.Message) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Message) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PullQuery provides a mock function with given fields: validatorID, requestID, containerID
func (_m *SnowmanEngine) PullQuery(validatorID ids.ShortID, requestID uint32, containerID ids.ID) error {
	ret := _m.Called(validatorID, requestID, containerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32, ids.ID) error); ok {
		r0 = rf(validatorID, requestID, containerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PushQuery provides a mock function with given fields: validatorID, requestID, containerID, container
func (_m *SnowmanEngine) PushQuery(validatorID ids.ShortID, requestID uint32, containerID ids.ID, container []byte) error {
	ret := _m.Called(validatorID, requestID, containerID, container)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32, ids.ID, []byte) error); ok {
		r0 = rf(validatorID, requestID, containerID, container)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Put provides a mock function with given fields: validatorID, requestID, containerID, container
func (_m *SnowmanEngine) Put(validatorID ids.ShortID, requestID uint32, containerID ids.ID, container []byte) error {
	ret := _m.Called(validatorID, requestID, containerID, container)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32, ids.ID, []byte) error); ok {
		r0 = rf(validatorID, requestID, containerID, container)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueryFailed provides a mock function with given fields: validatorID, requestID
func (_m *SnowmanEngine) QueryFailed(validatorID ids.ShortID, requestID uint32) error {
	ret := _m.Called(validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.ShortID, uint32) error); ok {
		r0 = rf(validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Shutdown provides a mock function with given fields:
func (_m *SnowmanEngine) Shutdown() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Startup provides a mock function with given fields:
func (_m *SnowmanEngine) Startup() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
