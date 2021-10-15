// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package router

import (
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/message"
	"github.com/ava-labs/avalanchego/utils/constants"
)

type messageWrap struct {
	inMsg     message.InboundMessage // Must always be set
	nodeID    ids.ShortID            // Must always be set
	requestID uint32
	received  time.Time // Time this message was received
	deadline  time.Time // Time this message must be responded to
}

func (m messageWrap) doneHandling() {
	if m.inMsg != nil {
		m.inMsg.OnFinishedHandling()
	}
}

// IsPeriodic returns true if this message is of a type that is sent on a
// periodic basis.
func (m messageWrap) IsPeriodic() bool {
	return m.requestID == constants.GossipMsgRequestID ||
		m.inMsg.Op() == message.GossipRequest
}
