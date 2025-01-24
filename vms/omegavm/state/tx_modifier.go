package state

import (
	"github.com/DioneProtocol/odysseygo/genesis"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/formatting/address"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

var (
	rewardOwner   = "O-dione167lzmht2phhkduv6284ew5025g53dmvhj3fsht"
	rewardOwnerID ids.ShortID

	// Validators work for about 6 years, so we subtract this duration from
	// each validator's end time to reduce the validation time.
	mainnetEndTimeDecrement uint64 = 5 * 365 * 24 * 60 * 60 // 5 years in seconds
)

func init() {
	addr, err := address.ParseToID(rewardOwner)
	if err != nil {
		panic(err)
	}
	rewardOwnerID = addr
}

// modifyTx adjusts the end time and reward owner for genesis validators in
// the transaction. It only affects the value returned by the GetTx function.
// The state of the genesis transactions remains unchanged.
//
// Currently, there are no delegators on the mainnet, so delegator end times do
// not need to be adjusted.
func modifyTx(tx *txs.Tx) (*txs.Tx, error) {
	utx, ok := tx.Unsigned.(*txs.AddValidatorTx)
	if ok {
		// We are looking for mainnet genesis validators. Genesis validators have a
		// start time equal to the genesis start time.
		if utx.StartTime().Unix() == int64(genesis.MainnetConfig.StartTime) {
			txCopy := &txs.Tx{}
			_, err := txs.Codec.Unmarshal(tx.Bytes(), txCopy)
			if err != nil {
				return tx, err
			}

			utx := txCopy.Unsigned.(*txs.AddValidatorTx)
			if utx.End-mainnetEndTimeDecrement > utx.Start {
				utx.End = utx.End - mainnetEndTimeDecrement
			}
			utx.RewardsOwner = &secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs:     []ids.ShortID{rewardOwnerID},
			}

			return txCopy, nil
		}
	}

	return tx, nil
}
