// (c) 2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package throttling

import (
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

// TODO add tests for max processing messages

func TestInboundMsgThrottler(t *testing.T) {
	assert := assert.New(t)
	config := InboundMsgThrottlerConfig{
		MaxProcessingMsgsPerNode: 1024,
		MsgByteThrottlerConfig: MsgByteThrottlerConfig{
			VdrAllocSize:        1024,
			AtLargeAllocSize:    1024,
			NodeMaxAtLargeBytes: 1024,
		},
	}
	vdrs := validators.NewSet()
	vdr1ID := ids.GenerateTestShortID()
	vdr2ID := ids.GenerateTestShortID()
	assert.NoError(vdrs.AddWeight(vdr1ID, 1))
	assert.NoError(vdrs.AddWeight(vdr2ID, 1))
	throttlerIntf, err := NewInboundMsgThrottler(
		&logging.Log{},
		"",
		prometheus.NewRegistry(),
		vdrs,
		config,
	)
	assert.NoError(err)

	// Make sure NewSybilInboundMsgThrottler works
	throttler := throttlerIntf.(*inboundMsgThrottler)
	assert.Equal(config.VdrAllocSize, throttler.byteThrottler.maxVdrBytes)
	assert.Equal(config.VdrAllocSize, throttler.byteThrottler.remainingVdrBytes)
	assert.Equal(config.AtLargeAllocSize, throttler.byteThrottler.remainingAtLargeBytes)
	assert.NotNil(throttler.byteThrottler.nodeToVdrBytesUsed)
	assert.NotNil(throttler.byteThrottler.log)
	assert.NotNil(throttler.byteThrottler.vdrs)
	assert.NotNil(throttler.byteThrottler.metrics)

	// Take from at-large allocation.
	// Should return immediately.
	throttlerIntf.Acquire(1, vdr1ID)
	assert.EqualValues(config.AtLargeAllocSize-1, throttler.byteThrottler.remainingAtLargeBytes)
	assert.EqualValues(config.VdrAllocSize, throttler.byteThrottler.remainingVdrBytes)
	assert.Len(throttler.byteThrottler.nodeToVdrBytesUsed, 0)
	assert.Len(throttler.byteThrottler.nodeToAtLargeBytesUsed, 1)
	assert.EqualValues(1, throttler.byteThrottler.nodeToAtLargeBytesUsed[vdr1ID])

	// Release the bytes
	throttlerIntf.Release(1, vdr1ID)
	assert.EqualValues(config.AtLargeAllocSize, throttler.byteThrottler.remainingAtLargeBytes)
	assert.EqualValues(config.VdrAllocSize, throttler.byteThrottler.remainingVdrBytes)
	assert.Len(throttler.byteThrottler.nodeToVdrBytesUsed, 0)
	assert.Len(throttler.byteThrottler.nodeToAtLargeBytesUsed, 0)

	// Use all the at-large allocation bytes and 1 of the validator allocation bytes
	// Should return immediately.
	throttlerIntf.Acquire(config.AtLargeAllocSize+1, vdr1ID)
	// vdr1 at-large bytes used: 1024. Validator bytes used: 1
	assert.EqualValues(0, throttler.byteThrottler.remainingAtLargeBytes)
	assert.EqualValues(config.VdrAllocSize-1, throttler.byteThrottler.remainingVdrBytes)
	assert.EqualValues(throttler.byteThrottler.nodeToVdrBytesUsed[vdr1ID], 1)
	assert.Len(throttler.byteThrottler.nodeToVdrBytesUsed, 1)
	assert.Len(throttler.byteThrottler.nodeToAtLargeBytesUsed, 1)
	assert.EqualValues(config.AtLargeAllocSize, throttler.byteThrottler.nodeToAtLargeBytesUsed[vdr1ID])

	// The other validator should be able to acquire half the validator allocation.
	// Should return immediately.
	throttlerIntf.Acquire(config.AtLargeAllocSize/2, vdr2ID)
	// vdr2 at-large bytes used: 0. Validator bytes used: 512
	assert.EqualValues(config.VdrAllocSize/2-1, throttler.byteThrottler.remainingVdrBytes)
	assert.EqualValues(throttler.byteThrottler.nodeToVdrBytesUsed[vdr1ID], 1)
	assert.EqualValues(throttler.byteThrottler.nodeToVdrBytesUsed[vdr2ID], config.VdrAllocSize/2)
	assert.Len(throttler.byteThrottler.nodeToVdrBytesUsed, 2)
	assert.Len(throttler.byteThrottler.nodeToAtLargeBytesUsed, 1)
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs, 0)
	assert.EqualValues(0, throttler.byteThrottler.waitingToAcquire.Len())

	// vdr1 should be able to acquire the rest of the validator allocation
	// Should return immediately.
	throttlerIntf.Acquire(config.VdrAllocSize/2-1, vdr1ID)
	// vdr1 at-large bytes used: 1024. Validator bytes used: 512
	assert.EqualValues(throttler.byteThrottler.nodeToVdrBytesUsed[vdr1ID], config.VdrAllocSize/2)
	assert.Len(throttler.byteThrottler.nodeToAtLargeBytesUsed, 1)
	assert.EqualValues(config.AtLargeAllocSize, throttler.byteThrottler.nodeToAtLargeBytesUsed[vdr1ID])

	// Trying to take more bytes for either node should block
	vdr1Done := make(chan struct{})
	go func() {
		throttlerIntf.Acquire(1, vdr1ID)
		vdr1Done <- struct{}{}
	}()
	select {
	case <-vdr1Done:
		t.Fatal("should block on acquiring any more bytes")
	case <-time.After(50 * time.Millisecond):
	}
	throttler.byteThrottler.lock.Lock()
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs, 1)
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs[vdr1ID], 1)
	assert.EqualValues(1, throttler.byteThrottler.waitingToAcquire.Len())
	_, exists := throttler.byteThrottler.waitingToAcquire.Get(throttler.byteThrottler.nodeToWaitingMsgIDs[vdr1ID][0])
	assert.True(exists)
	throttler.byteThrottler.lock.Unlock()

	vdr2Done := make(chan struct{})
	go func() {
		throttlerIntf.Acquire(1, vdr2ID)
		vdr2Done <- struct{}{}
	}()
	select {
	case <-vdr2Done:
		t.Fatal("should block on acquiring any more bytes")
	case <-time.After(50 * time.Millisecond):
	}
	throttler.byteThrottler.lock.Lock()
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs, 2)
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs[vdr2ID], 1)
	assert.EqualValues(2, throttler.byteThrottler.waitingToAcquire.Len())
	_, exists = throttler.byteThrottler.waitingToAcquire.Get(throttler.byteThrottler.nodeToWaitingMsgIDs[vdr2ID][0])
	assert.True(exists)
	throttler.byteThrottler.lock.Unlock()

	nonVdrID := ids.GenerateTestShortID()
	nonVdrDone := make(chan struct{})
	go func() {
		throttlerIntf.Acquire(1, nonVdrID)
		nonVdrDone <- struct{}{}
	}()
	select {
	case <-nonVdrDone:
		t.Fatal("should block on acquiring any more bytes")
	case <-time.After(50 * time.Millisecond):
	}
	throttler.byteThrottler.lock.Lock()
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs, 3)
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs[nonVdrID], 1)
	assert.EqualValues(3, throttler.byteThrottler.waitingToAcquire.Len())
	_, exists = throttler.byteThrottler.waitingToAcquire.Get(throttler.byteThrottler.nodeToWaitingMsgIDs[nonVdrID][0])
	assert.True(exists)
	throttler.byteThrottler.lock.Unlock()

	// Release config.MaxAtLargeBytes+1 bytes
	// When the choice exists, bytes should be given back to the validator allocation
	// rather than the at-large allocation.
	throttlerIntf.Release(config.AtLargeAllocSize+1, vdr1ID)

	// The Acquires that blocked above should have returned
	<-vdr1Done
	<-vdr2Done
	<-nonVdrDone

	assert.EqualValues(config.NodeMaxAtLargeBytes/2, throttler.byteThrottler.remainingVdrBytes)
	assert.Len(throttler.byteThrottler.nodeToAtLargeBytesUsed, 3) // vdr1, vdr2, nonVdrID
	assert.EqualValues(config.AtLargeAllocSize/2, throttler.byteThrottler.nodeToAtLargeBytesUsed[vdr1ID])
	assert.EqualValues(1, throttler.byteThrottler.nodeToAtLargeBytesUsed[vdr2ID])
	assert.EqualValues(1, throttler.byteThrottler.nodeToAtLargeBytesUsed[nonVdrID])
	assert.Len(throttler.byteThrottler.nodeToVdrBytesUsed, 1)
	assert.EqualValues(0, throttler.byteThrottler.nodeToVdrBytesUsed[vdr1ID])
	assert.EqualValues(config.AtLargeAllocSize/2-2, throttler.byteThrottler.remainingAtLargeBytes)
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs, 0)
	assert.EqualValues(0, throttler.byteThrottler.waitingToAcquire.Len())

	// Non-validator should be able to take the rest of the at-large bytes
	throttlerIntf.Acquire(config.AtLargeAllocSize/2-2, nonVdrID)
	assert.EqualValues(0, throttler.byteThrottler.remainingAtLargeBytes)
	assert.EqualValues(config.AtLargeAllocSize/2-1, throttler.byteThrottler.nodeToAtLargeBytesUsed[nonVdrID])
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs, 0)
	assert.EqualValues(0, throttler.byteThrottler.waitingToAcquire.Len())

	// But should block on subsequent Acquires
	go func() {
		throttlerIntf.Acquire(1, nonVdrID)
		nonVdrDone <- struct{}{}
	}()
	select {
	case <-nonVdrDone:
		t.Fatal("should block on acquiring any more bytes")
	case <-time.After(50 * time.Millisecond):
	}
	throttler.byteThrottler.lock.Lock()
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs, 1)
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs[nonVdrID], 1)
	assert.EqualValues(1, throttler.byteThrottler.waitingToAcquire.Len())
	_, exists = throttler.byteThrottler.waitingToAcquire.Get(throttler.byteThrottler.nodeToWaitingMsgIDs[nonVdrID][0])
	assert.True(exists)
	throttler.byteThrottler.lock.Unlock()

	// Release all of vdr2's messages
	throttlerIntf.Release(config.AtLargeAllocSize/2, vdr2ID)
	throttlerIntf.Release(1, vdr2ID)

	<-nonVdrDone

	assert.EqualValues(0, throttler.byteThrottler.nodeToAtLargeBytesUsed[vdr2ID])
	assert.EqualValues(config.VdrAllocSize, throttler.byteThrottler.remainingVdrBytes)
	assert.Len(throttler.byteThrottler.nodeToVdrBytesUsed, 0)
	assert.EqualValues(0, throttler.byteThrottler.remainingAtLargeBytes)
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs, 0)
	assert.EqualValues(0, throttler.byteThrottler.waitingToAcquire.Len())

	// Release all of vdr1's messages
	throttlerIntf.Release(1, vdr1ID)
	throttlerIntf.Release(config.AtLargeAllocSize/2-1, vdr1ID)
	assert.Len(throttler.byteThrottler.nodeToVdrBytesUsed, 0)
	assert.EqualValues(config.VdrAllocSize, throttler.byteThrottler.remainingVdrBytes)
	assert.EqualValues(config.AtLargeAllocSize/2, throttler.byteThrottler.remainingAtLargeBytes)
	assert.EqualValues(0, throttler.byteThrottler.nodeToAtLargeBytesUsed[vdr1ID])
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs, 0)
	assert.EqualValues(0, throttler.byteThrottler.waitingToAcquire.Len())

	// Release nonVdr's messages
	throttlerIntf.Release(1, nonVdrID)
	throttlerIntf.Release(1, nonVdrID)
	throttlerIntf.Release(config.AtLargeAllocSize/2-2, nonVdrID)
	assert.Len(throttler.byteThrottler.nodeToVdrBytesUsed, 0)
	assert.EqualValues(config.VdrAllocSize, throttler.byteThrottler.remainingVdrBytes)
	assert.EqualValues(config.AtLargeAllocSize, throttler.byteThrottler.remainingAtLargeBytes)
	assert.Len(throttler.byteThrottler.nodeToAtLargeBytesUsed, 0)
	assert.EqualValues(0, throttler.byteThrottler.nodeToAtLargeBytesUsed[nonVdrID])
	assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs, 0)
	assert.EqualValues(0, throttler.byteThrottler.waitingToAcquire.Len())
}

// Ensure that the limit on taking from the at-large allocation is enforced
func TestSybilMsgThrottlerMaxNonVdr(t *testing.T) {
	assert := assert.New(t)
	config := InboundMsgThrottlerConfig{
		MaxProcessingMsgsPerNode: 1024,
		MsgByteThrottlerConfig: MsgByteThrottlerConfig{
			VdrAllocSize:        100,
			AtLargeAllocSize:    100,
			NodeMaxAtLargeBytes: 10,
		},
	}
	vdrs := validators.NewSet()
	vdr1ID := ids.GenerateTestShortID()
	assert.NoError(vdrs.AddWeight(vdr1ID, 1))
	throttlerIntf, err := NewInboundMsgThrottler(
		&logging.Log{},
		"",
		prometheus.NewRegistry(),
		vdrs,
		config,
	)
	assert.NoError(err)
	throttler := throttlerIntf.(*inboundMsgThrottler)
	nonVdrNodeID1 := ids.GenerateTestShortID()
	throttlerIntf.Acquire(config.NodeMaxAtLargeBytes, nonVdrNodeID1)

	// Acquiring more should block
	nonVdrDone := make(chan struct{})
	go func() {
		throttlerIntf.Acquire(1, nonVdrNodeID1)
		nonVdrDone <- struct{}{}
	}()
	select {
	case <-nonVdrDone:
		t.Fatal("should block on acquiring any more bytes")
	case <-time.After(50 * time.Millisecond):
	}

	// A different non-validator should be able to acquire
	nonVdrNodeID2 := ids.GenerateTestShortID()
	throttlerIntf.Acquire(config.NodeMaxAtLargeBytes, nonVdrNodeID2)

	// Acquiring more should block
	go func() {
		throttlerIntf.Acquire(1, nonVdrNodeID1)
		nonVdrDone <- struct{}{}
	}()
	select {
	case <-nonVdrDone:
		t.Fatal("should block on acquiring any more bytes")
	case <-time.After(50 * time.Millisecond):
	}

	// Validator should only be able to take [MaxAtLargeBytes]
	throttlerIntf.Acquire(config.NodeMaxAtLargeBytes+1, vdr1ID)
	assert.EqualValues(config.NodeMaxAtLargeBytes, throttler.byteThrottler.nodeToAtLargeBytesUsed[vdr1ID])
	assert.EqualValues(1, throttler.byteThrottler.nodeToVdrBytesUsed[vdr1ID])
	assert.EqualValues(config.NodeMaxAtLargeBytes, throttler.byteThrottler.nodeToAtLargeBytesUsed[nonVdrNodeID1])
	assert.EqualValues(config.NodeMaxAtLargeBytes, throttler.byteThrottler.nodeToAtLargeBytesUsed[nonVdrNodeID2])
	assert.EqualValues(config.AtLargeAllocSize-config.NodeMaxAtLargeBytes*3, throttler.byteThrottler.remainingAtLargeBytes)
}

// Test that messages waiting to be acquired by a given node
// are handled in FIFO order
func TestSybilMsgThrottlerFIFO(t *testing.T) {
	assert := assert.New(t)
	config := InboundMsgThrottlerConfig{
		MaxProcessingMsgsPerNode: 1024,
		MsgByteThrottlerConfig: MsgByteThrottlerConfig{
			VdrAllocSize:        1024,
			AtLargeAllocSize:    1024,
			NodeMaxAtLargeBytes: 1024,
		},
	}
	vdrs := validators.NewSet()
	vdr1ID := ids.GenerateTestShortID()
	assert.NoError(vdrs.AddWeight(vdr1ID, 1))
	nonVdrNodeID := ids.GenerateTestShortID()

	maxVdrBytes := config.VdrAllocSize + config.AtLargeAllocSize
	maxNonVdrBytes := config.AtLargeAllocSize
	// Test for both validator and non-validator
	for _, nodeID := range []ids.ShortID{vdr1ID, nonVdrNodeID} {
		maxBytes := maxVdrBytes
		if nodeID == nonVdrNodeID {
			maxBytes = maxNonVdrBytes
		}
		throttlerIntf, err := NewInboundMsgThrottler(
			&logging.Log{},
			"",
			prometheus.NewRegistry(),
			vdrs,
			config,
		)
		assert.NoError(err)
		throttler := throttlerIntf.(*inboundMsgThrottler)
		// node uses up all but 1 byte
		throttler.Acquire(maxBytes-1, nodeID)
		// node uses the last byte
		throttler.Acquire(1, nodeID)

		// First message wants to acquire a lot of bytes
		done := make(chan struct{})
		go func() {
			throttler.Acquire(maxBytes-1, nodeID)
			done <- struct{}{}
		}()
		select {
		case <-done:
			t.Fatal("should block on acquiring any more bytes")
		case <-time.After(50 * time.Millisecond):
		}

		// Next message only wants to acquire 1 byte
		go func() {
			throttler.Acquire(1, nodeID)
			done <- struct{}{}
		}()
		select {
		case <-done:
			t.Fatal("should block on acquiring any more bytes")
		case <-time.After(50 * time.Millisecond):
		}

		// Release 1 byte
		throttler.Release(1, nodeID)
		// Byte should have gone toward first message
		assert.EqualValues(2, throttler.byteThrottler.waitingToAcquire.Len())
		assert.Len(throttler.byteThrottler.nodeToWaitingMsgIDs[nodeID], 2)
		firstMsgID := throttler.byteThrottler.nodeToWaitingMsgIDs[nodeID][0]
		firstMsg, exists := throttler.byteThrottler.waitingToAcquire.Get(firstMsgID)
		assert.True(exists)
		assert.EqualValues(maxBytes-2, firstMsg.(*msgMetadata).bytesNeeded)

		// Since messages are processed FIFO for a given validator,
		// the first message should return from Acquire first
		select {
		case <-done:
			t.Fatal("should still be blocking")
		case <-time.After(50 * time.Millisecond):
		}

		// Release the rest of the bytes
		throttler.Release(maxBytes-1, nodeID)
		// Both should be done acquiring now
		<-done
		<-done
	}
}
