// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package p

import (
	stdcontext "context"

	"github.com/DioneProtocol/odysseygo/api/info"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/vms/alpha"
)

var _ Context = (*context)(nil)

type Context interface {
	NetworkID() uint32
	DIONEAssetID() ids.ID
	BaseTxFee() uint64
	CreateSubnetTxFee() uint64
	TransformSubnetTxFee() uint64
	CreateBlockchainTxFee() uint64
	AddPrimaryNetworkValidatorFee() uint64
	AddSubnetValidatorFee() uint64
}

type context struct {
	networkID                     uint32
	dioneAssetID                  ids.ID
	baseTxFee                     uint64
	createSubnetTxFee             uint64
	transformSubnetTxFee          uint64
	createBlockchainTxFee         uint64
	addPrimaryNetworkValidatorFee uint64
	addSubnetValidatorFee         uint64
}

func NewContextFromURI(ctx stdcontext.Context, uri string) (Context, error) {
	infoClient := info.NewClient(uri)
	aChainClient := alpha.NewClient(uri, "A")
	return NewContextFromClients(ctx, infoClient, aChainClient)
}

func NewContextFromClients(
	ctx stdcontext.Context,
	infoClient info.Client,
	aChainClient alpha.Client,
) (Context, error) {
	networkID, err := infoClient.GetNetworkID(ctx)
	if err != nil {
		return nil, err
	}

	asset, err := aChainClient.GetAssetDescription(ctx, "DIONE")
	if err != nil {
		return nil, err
	}

	txFees, err := infoClient.GetTxFee(ctx)
	if err != nil {
		return nil, err
	}

	return NewContext(
		networkID,
		asset.AssetID,
		uint64(txFees.TxFee),
		uint64(txFees.CreateSubnetTxFee),
		uint64(txFees.TransformSubnetTxFee),
		uint64(txFees.CreateBlockchainTxFee),
		uint64(txFees.AddPrimaryNetworkValidatorFee),
		uint64(txFees.AddSubnetValidatorFee),
	), nil
}

func NewContext(
	networkID uint32,
	dioneAssetID ids.ID,
	baseTxFee uint64,
	createSubnetTxFee uint64,
	transformSubnetTxFee uint64,
	createBlockchainTxFee uint64,
	addPrimaryNetworkValidatorFee uint64,
	addSubnetValidatorFee uint64,
) Context {
	return &context{
		networkID:                     networkID,
		dioneAssetID:                  dioneAssetID,
		baseTxFee:                     baseTxFee,
		createSubnetTxFee:             createSubnetTxFee,
		transformSubnetTxFee:          transformSubnetTxFee,
		createBlockchainTxFee:         createBlockchainTxFee,
		addPrimaryNetworkValidatorFee: addPrimaryNetworkValidatorFee,
		addSubnetValidatorFee:         addSubnetValidatorFee,
	}
}

func (c *context) NetworkID() uint32 {
	return c.networkID
}

func (c *context) DIONEAssetID() ids.ID {
	return c.dioneAssetID
}

func (c *context) BaseTxFee() uint64 {
	return c.baseTxFee
}

func (c *context) CreateSubnetTxFee() uint64 {
	return c.createSubnetTxFee
}

func (c *context) TransformSubnetTxFee() uint64 {
	return c.transformSubnetTxFee
}

func (c *context) CreateBlockchainTxFee() uint64 {
	return c.createBlockchainTxFee
}

func (c *context) AddPrimaryNetworkValidatorFee() uint64 {
	return c.addPrimaryNetworkValidatorFee
}

func (c *context) AddSubnetValidatorFee() uint64 {
	return c.addSubnetValidatorFee
}
