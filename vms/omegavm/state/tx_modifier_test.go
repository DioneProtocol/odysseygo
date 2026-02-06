package state

import (
	"testing"

	nodegenesis "github.com/DioneProtocol/odysseygo/genesis"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	blocks "github.com/DioneProtocol/odysseygo/vms/omegavm/blocks"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/genesis"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/reward"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/stretchr/testify/require"
)

func TestTxModifier(t *testing.T) {
	require := require.New(t)
	newState, db := newUninitializedState(require)

	endTime := nodegenesis.MainnetConfig.StartTime + nodegenesis.MainnetConfig.InitialStakeDuration
	genesisNodeId := ids.GenerateTestNodeID()
	genesisValidator := &txs.AddValidatorTx{
		Validator: txs.Validator{
			NodeID: genesisNodeId,
			Start:  uint64(nodegenesis.MainnetConfig.StartTime),
			End:    uint64(endTime),
			Wght:   units.Dione,
		},
		StakeOuts: []*dione.TransferableOutput{
			{
				Asset: dione.Asset{ID: initialTxID},
				Out: &secp256k1fx.TransferOutput{
					Amt: units.Dione,
				},
			},
		},
		RewardsOwner: &secp256k1fx.OutputOwners{
			Addrs:     []ids.ShortID{rewardOwnerID},
			Threshold: uint32(1),
		},
		DelegationShares: reward.PercentDenominator,
	}
	genesisValidatorTx := &txs.Tx{Unsigned: genesisValidator}
	require.NoError(genesisValidatorTx.Initialize(txs.Codec))

	anotherValidator := &txs.AddValidatorTx{
		Validator: txs.Validator{
			NodeID: initialNodeID,
			Start:  uint64(initialTime.Unix()),
			End:    uint64(initialValidatorEndTime.Unix()),
			Wght:   units.Dione,
		},
		StakeOuts: []*dione.TransferableOutput{
			{
				Asset: dione.Asset{ID: initialTxID},
				Out: &secp256k1fx.TransferOutput{
					Amt: units.Dione,
				},
			},
		},
		RewardsOwner:     &secp256k1fx.OutputOwners{},
		DelegationShares: reward.PercentDenominator,
	}
	anotherValidatorTx := &txs.Tx{Unsigned: anotherValidator}
	require.NoError(anotherValidatorTx.Initialize(txs.Codec))

	genesisState := &genesis.State{
		UTXOs: []*dione.UTXO{
			{
				UTXOID: dione.UTXOID{
					TxID:        initialTxID,
					OutputIndex: 0,
				},
				Asset: dione.Asset{ID: initialTxID},
				Out: &secp256k1fx.TransferOutput{
					Amt: units.Schmeckle,
				},
			},
		},
		Validators: []*txs.Tx{
			genesisValidatorTx,
			anotherValidatorTx,
		},
		Chains:        []*txs.Tx{},
		Timestamp:     uint64(initialTime.Unix()),
		InitialSupply: units.Schmeckle + units.Dione,
	}
	genesisBlkID := ids.GenerateTestID()
	genesisBlk, err := blocks.NewApricotCommitBlock(genesisBlkID, 0)
	require.NoError(err)
	require.NoError(newState.(*state).syncGenesis(genesisBlk, genesisState))
	require.NoError(newState.(*state).doneInit())
	require.NoError(newState.(*state).Commit())

	loadedState := newStateFromDB(require, db)
	require.NoError(loadedState.(*state).load())

	for _, s := range []State{newState, loadedState} {
		tx, _, err := s.GetTx(genesisValidatorTx.TxID)
		require.NoError(err)

		currentValidators, ok := s.(*state).currentStakers.validators[constants.PrimaryNetworkID]
		require.True(ok)

		utx := tx.Unsigned.(*txs.AddValidatorTx)
		owner := utx.RewardsOwner.(*secp256k1fx.OutputOwners)
		require.Equal(genesisValidator.Validator.NodeID, utx.NodeID())
		require.Equal(genesisValidator.Validator.Start, utx.Start)
		require.Equal(genesisValidator.Validator.End, utx.End)
		require.Equal(genesisValidator.Validator.Wght, utx.Wght)
		require.Equal(len(owner.Addrs), 1)
		require.Equal(owner.Addrs[0], rewardOwnerID)
		require.Equal(owner.Threshold, uint32(1))
		require.Equal(owner.Locktime, uint64(0))

		validator, err := s.GetCurrentValidator(constants.PrimaryNetworkID, genesisNodeId)
		require.NoError(err)

		require.Equal(validator.TxID, genesisValidatorTx.TxID)
		require.Equal(validator.NodeID, genesisNodeId)
		require.Equal(validator.Weight, utx.Wght)
		require.Equal(validator.StartTime.Unix(), int64(utx.Start))
		require.Equal(validator.EndTime.Unix(), int64(utx.End))
		require.Equal(validator.NextTime.Unix(), int64(utx.End))

		genesisValidatorBaseStaker, ok := currentValidators[genesisNodeId]
		require.True(ok)
		require.Equal(genesisValidatorBaseStaker.validator.EndTime.Unix(), int64(utx.End))

		tx, _, err = s.GetTx(anotherValidatorTx.TxID)
		require.NoError(err)

		utx = tx.Unsigned.(*txs.AddValidatorTx)
		owner = utx.RewardsOwner.(*secp256k1fx.OutputOwners)
		require.Equal(anotherValidator.Validator.NodeID, utx.NodeID())
		require.Equal(anotherValidator.Validator.Start, utx.Start)
		require.Equal(anotherValidator.Validator.End, utx.End)
		require.Equal(anotherValidator.Validator.Wght, utx.Wght)
		require.Equal(len(owner.Addrs), 0)
		require.Equal(owner.Threshold, uint32(0))
		require.Equal(owner.Locktime, uint64(0))

		validator, err = s.GetCurrentValidator(constants.PrimaryNetworkID, initialNodeID)
		require.NoError(err)

		require.Equal(validator.TxID, anotherValidatorTx.TxID)
		require.Equal(validator.NodeID, initialNodeID)
		require.Equal(validator.Weight, utx.Wght)
		require.Equal(validator.StartTime.Unix(), int64(utx.Start))
		require.Equal(validator.EndTime.Unix(), int64(utx.End))
		require.Equal(validator.NextTime.Unix(), int64(utx.End))

		anotherValidatorBaseStaker, ok := currentValidators[initialNodeID]
		require.True(ok)
		require.Equal(anotherValidatorBaseStaker.validator.EndTime.Unix(), int64(utx.End))
	}
}
