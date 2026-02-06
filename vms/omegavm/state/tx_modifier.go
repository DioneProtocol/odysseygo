package state

import (
	"time"

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

	// update time for delegator which has start time equal to genesis start time
	testnetDelegatorEndTime = time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)

	nonActiveValidatorStartTime = time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC)

	// update time for non active validator which has start time before nonActiveValidatorStartTime
	testnetNonActiveValidatorEndTimeDecrement uint64 = 1 * 365 * 24 * 60 * 60 // 1 year in seconds

	// update time for validator which has end time before testnetEndTimeDecrement
	testnetEndTimeDecrement uint64 = 5 * 365 * 24 * 60 * 60 // 5 years in seconds
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
// Currently, there are no delegators on the testnet, so delegator end times do
// not need to be adjusted.
func (s *state) modifyTx(tx *txs.Tx) (*txs.Tx, error) {
	// update end time when pyruni is activated
	if !s.cfg.IsPyruniActivated(s.timestamp) {
		return tx, nil
	}

	activationTime := s.cfg.PyruniTime.Unix()
	switch utx := tx.Unsigned.(type) {
	case *txs.AddValidatorTx:
		// We are looking for testnet genesis validators. Genesis validators have a
		// start time equal to the genesis start time.
		if utx.StartTime().Unix() == int64(genesis.TestnetConfig.StartTime) {
			return reduceValidatorEndTime(tx, testnetEndTimeDecrement)
		}

		if (utx.StartTime().Unix() <= nonActiveValidatorStartTime.Unix() && utx.EndTime().Unix() > activationTime) && utx.StartTime().Unix() != int64(genesis.TestnetConfig.StartTime) {
			return reduceValidatorEndTime(tx, testnetNonActiveValidatorEndTimeDecrement)
		}
	case *txs.AddDelegatorTx:
		vdrStartTime, err := s.GetStartTime(utx.NodeID(), utx.SubnetID())
		if err != nil {
			return tx, nil
		}
		if checkReductionRequiredForDelegator(vdrStartTime.Unix(), utx.StartTime().Unix(), utx.EndTime().Unix(), activationTime) {
			txCopy := &txs.Tx{}
			_, err := txs.Codec.Unmarshal(tx.Bytes(), txCopy)
			if err != nil {
				return tx, err
			}
	
			utx := txCopy.Unsigned.(*txs.AddDelegatorTx)
			utx.End = uint64(testnetDelegatorEndTime.Unix())
			return txCopy, nil
	
		}
	case *txs.AddPermissionlessDelegatorTx:
		vdrStartTime, err := s.GetStartTime(utx.NodeID(), utx.SubnetID())
		if err != nil {
			return tx, nil
		}
		if checkReductionRequiredForDelegator(vdrStartTime.Unix(), utx.StartTime().Unix(), utx.EndTime().Unix(), activationTime) {
			txCopy := &txs.Tx{}
			_, err := txs.Codec.Unmarshal(tx.Bytes(), txCopy)
			if err != nil {
				return tx, err
			}
	
			utx := txCopy.Unsigned.(*txs.AddPermissionlessDelegatorTx)
			utx.End = uint64(testnetDelegatorEndTime.Unix())
	
			return txCopy, nil
		}
	}

	return tx, nil
}

func reduceValidatorEndTime(tx *txs.Tx, endTimeDecrement uint64) (*txs.Tx, error) {
	txCopy := &txs.Tx{}
	_, err := txs.Codec.Unmarshal(tx.Bytes(), txCopy)
	if err != nil {
		return tx, err
	}

	utx := txCopy.Unsigned.(*txs.AddValidatorTx)
	if int64(utx.End-endTimeDecrement) > int64(utx.Start) {
		utx.End = utx.End - endTimeDecrement
	}
	utx.RewardsOwner = &secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs:     []ids.ShortID{rewardOwnerID},
	}

	return txCopy, nil
}

func checkReductionRequiredForDelegator(vdrStartTime int64, startTime int64, endTime int64, activationTime int64) bool {
	return vdrStartTime == int64(genesis.TestnetConfig.StartTime) && (startTime < activationTime && endTime > activationTime)
}
