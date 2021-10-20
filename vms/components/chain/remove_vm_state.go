package chain

import (
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/wrappers"
)

// In this file methods specific of RemoteVM are implemented

func (s *State) GetAncestors(
	blkID ids.ID,
	maxBlocksNum int,
	maxBlocksSize int,
	maxBlocksRetrivalTime time.Duration,
) [][]byte {
	res := make([][]byte, 0, maxBlocksNum)
	currentByteLength := 0
	startTime := time.Now()
	var currentDuration time.Duration

	for cnt := 0; ; cnt++ {
		currentDuration = time.Since(startTime)
		if cnt >= maxBlocksNum || currentDuration >= maxBlocksRetrivalTime {
			return res // return what we have
		}

		blk, ok := s.getCachedBlock(blkID)
		if !ok {
			if _, ok := s.missingBlocks.Get(blkID); ok {
				// currently requested block is not cached nor available over wire.
				return res // return what we have
			}

			// currently requested block and its ancestors may be available over wire
			break // done with cache, move to wire
		}

		blkBytes := blk.Bytes()
		// Ensure response size isn't too large. Include wrappers.IntLen because the size of the message
		// is included with each container, and the size is repr. by an int.
		if newLen := currentByteLength + wrappers.IntLen + len(blkBytes); newLen < maxBlocksSize {
			res = append(res, blkBytes)
			currentByteLength = newLen
			blkID = blk.Parent()
			continue
		}

		// reached maximum response size, return what we have
		return res
	}

	wireBlkID := blkID
	wireMaxBlkNum := maxBlocksNum - len(res)
	wireMaxBlocksSize := maxBlocksSize - currentByteLength
	wireMaxBlocksRetrivalTime := time.Duration(maxBlocksRetrivalTime.Nanoseconds() - currentDuration.Nanoseconds())

	wireRes := s.getAncestors(wireBlkID, wireMaxBlkNum, wireMaxBlocksSize, wireMaxBlocksRetrivalTime)
	res = append(res, wireRes...)
	return res
}
