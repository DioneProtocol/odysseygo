// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package network

import (
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils"
)

// Builder extends a Codec to build messages safely
type Builder struct {
	Codec
	// [getByteSlice] must not be nil.
	// [getByteSlice] may return nil.
	// [getByteSlice] must be safe for concurrent access by multiple goroutines.
	getByteSlice func() []byte
}

// GetVersion message
func (m Builder) GetVersion() (Msg, error) {
	buf := m.getByteSlice()
	return m.Pack(buf, GetVersion, nil, false)
}

// Version message
func (m Builder) Version(networkID, nodeID uint32, myTime uint64, ip utils.IPDesc, myVersion string, myVersionTime uint64, sig []byte) (Msg, error) {
	buf := m.getByteSlice()
	return m.Pack(
		buf,
		Version,
		map[Field]interface{}{
			NetworkID:   networkID,
			NodeID:      nodeID,
			MyTime:      myTime,
			IP:          ip,
			VersionStr:  myVersion,
			VersionTime: myVersionTime,
			SigBytes:    sig,
		},
		false,
	)
}

// GetPeerList message
func (m Builder) GetPeerList() (Msg, error) {
	buf := m.getByteSlice()
	return m.Pack(buf, GetPeerList, nil, false)
}

func (m Builder) PeerList(peers []utils.IPCertDesc) (Msg, error) {
	buf := m.getByteSlice()
	return m.Pack(
		buf,
		PeerList,
		map[Field]interface{}{
			SignedPeers: peers,
		},
		false,
	)
}

// Ping message
func (m Builder) Ping() (Msg, error) {
	buf := m.getByteSlice()
	return m.Pack(buf, Ping, nil, false)
}

// Pong message
func (m Builder) Pong() (Msg, error) {
	buf := m.getByteSlice()
	return m.Pack(buf, Pong, nil, false)
}

// GetAcceptedFrontier message
func (m Builder) GetAcceptedFrontier(chainID ids.ID, requestID uint32, deadline uint64) (Msg, error) {
	buf := m.getByteSlice()
	return m.Pack(
		buf,
		GetAcceptedFrontier,
		map[Field]interface{}{
			ChainID:   chainID[:],
			RequestID: requestID,
			Deadline:  deadline,
		},
		false,
	)
}

// AcceptedFrontier message
func (m Builder) AcceptedFrontier(chainID ids.ID, requestID uint32, containerIDs []ids.ID) (Msg, error) {
	containerIDBytes := make([][]byte, len(containerIDs))
	for i, containerID := range containerIDs {
		copy := containerID
		containerIDBytes[i] = copy[:]
	}
	buf := m.getByteSlice()
	return m.Pack(
		buf,
		AcceptedFrontier,
		map[Field]interface{}{
			ChainID:      chainID[:],
			RequestID:    requestID,
			ContainerIDs: containerIDBytes,
		},
		false,
	)
}

// GetAccepted message
func (m Builder) GetAccepted(chainID ids.ID, requestID uint32, deadline uint64, containerIDs []ids.ID) (Msg, error) {
	containerIDBytes := make([][]byte, len(containerIDs))
	for i, containerID := range containerIDs {
		copy := containerID
		containerIDBytes[i] = copy[:]
	}
	buf := m.getByteSlice()
	return m.Pack(
		buf,
		GetAccepted,
		map[Field]interface{}{
			ChainID:      chainID[:],
			RequestID:    requestID,
			Deadline:     deadline,
			ContainerIDs: containerIDBytes,
		},
		false,
	)
}

// Accepted message
func (m Builder) Accepted(chainID ids.ID, requestID uint32, containerIDs []ids.ID) (Msg, error) {
	containerIDBytes := make([][]byte, len(containerIDs))
	for i, containerID := range containerIDs {
		copy := containerID
		containerIDBytes[i] = copy[:]
	}
	buf := m.getByteSlice()
	return m.Pack(
		buf,
		Accepted,
		map[Field]interface{}{
			ChainID:      chainID[:],
			RequestID:    requestID,
			ContainerIDs: containerIDBytes,
		},
		false,
	)
}

// GetAncestors message
func (m Builder) GetAncestors(chainID ids.ID, requestID uint32, deadline uint64, containerID ids.ID) (Msg, error) {
	buf := m.getByteSlice()
	return m.Pack(
		buf,
		GetAncestors,
		map[Field]interface{}{
			ChainID:     chainID[:],
			RequestID:   requestID,
			Deadline:    deadline,
			ContainerID: containerID[:],
		},
		false,
	)
}

// MultiPut message
func (m Builder) MultiPut(chainID ids.ID, requestID uint32, containers [][]byte) (Msg, error) {
	buf := m.getByteSlice()
	return m.Pack(
		buf,
		MultiPut,
		map[Field]interface{}{
			ChainID:             chainID[:],
			RequestID:           requestID,
			MultiContainerBytes: containers,
		},
		false,
	)
}

// Get message
func (m Builder) Get(chainID ids.ID, requestID uint32, deadline uint64, containerID ids.ID) (Msg, error) {
	buf := m.getByteSlice()
	return m.Pack(
		buf,
		Get,
		map[Field]interface{}{
			ChainID:     chainID[:],
			RequestID:   requestID,
			Deadline:    deadline,
			ContainerID: containerID[:],
		},
		false,
	)
}

// Put message
func (m Builder) Put(chainID ids.ID, requestID uint32, containerID ids.ID, container []byte) (Msg, error) {
	buf := m.getByteSlice()
	return m.Pack(
		buf,
		Put,
		map[Field]interface{}{
			ChainID:        chainID[:],
			RequestID:      requestID,
			ContainerID:    containerID[:],
			ContainerBytes: container,
		},
		false,
	)
}

// PushQuery message
func (m Builder) PushQuery(chainID ids.ID, requestID uint32, deadline uint64, containerID ids.ID, container []byte) (Msg, error) {
	buf := m.getByteSlice()
	return m.Pack(
		buf,
		PushQuery,
		map[Field]interface{}{
			ChainID:        chainID[:],
			RequestID:      requestID,
			Deadline:       deadline,
			ContainerID:    containerID[:],
			ContainerBytes: container,
		},
		false,
	)
}

// PullQuery message
func (m Builder) PullQuery(chainID ids.ID, requestID uint32, deadline uint64, containerID ids.ID) (Msg, error) {
	buf := m.getByteSlice()
	return m.Pack(
		buf,
		PullQuery,
		map[Field]interface{}{
			ChainID:     chainID[:],
			RequestID:   requestID,
			Deadline:    deadline,
			ContainerID: containerID[:],
		},
		false,
	)
}

// Chits message
func (m Builder) Chits(chainID ids.ID, requestID uint32, containerIDs []ids.ID) (Msg, error) {
	containerIDBytes := make([][]byte, len(containerIDs))
	for i, containerID := range containerIDs {
		copy := containerID
		containerIDBytes[i] = copy[:]
	}
	buf := m.getByteSlice()
	return m.Pack(
		buf,
		Chits,
		map[Field]interface{}{
			ChainID:      chainID[:],
			RequestID:    requestID,
			ContainerIDs: containerIDBytes,
		},
		false,
	)
}
