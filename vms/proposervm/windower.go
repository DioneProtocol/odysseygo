package proposervm

import (
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/validators"
)

// windower interfaces with P-Chain and it is responsible for:
// retrieving current P-Chain height
// calculate the start time for the block submission window of a given validator

type windower struct {
	validators.VM
	mockedValPos uint
}

func (w *windower) pChainHeight() (uint64, error) {
	return w.VM.GetCurrentHeight()
}

func (w *windower) BlkSubmissionDelay(pChainHeight uint64, valID ids.ShortID) time.Duration {
	// TODO:
	//       pick validators population at given pChainHeight
	//       if valID not in validator set, valPos = len(validators population)
	//       else pick random permutation, seed by pChainHeight???
	//       valPos is valID position in the permutation
	return time.Duration(w.mockedValPos) * BlkSubmissionWinLength
}
