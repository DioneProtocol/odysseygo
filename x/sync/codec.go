// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package sync

import (
	"github.com/dioneprotocol/dionego/codec"
	"github.com/dioneprotocol/dionego/codec/linearcodec"
	"github.com/dioneprotocol/dionego/utils/units"
	"github.com/dioneprotocol/dionego/utils/wrappers"
)

const (
	Version        = 0 // TODO danlaine unexport this
	maxMessageSize = 1 * units.MiB
)

var syncCodec codec.Manager

func init() {
	syncCodec = codec.NewManager(maxMessageSize)
	c := linearcodec.NewDefault()

	errs := wrappers.Errs{}
	errs.Add(
		c.RegisterType(&ChangeProofRequest{}),
		c.RegisterType(&RangeProofRequest{}),
		syncCodec.RegisterCodec(Version, c),
	)

	if errs.Errored() {
		panic(errs.Err)
	}
}
