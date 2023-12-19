// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package blocks

type Visitor interface {
	BanffAbortBlock(*BanffAbortBlock) error
	BanffCommitBlock(*BanffCommitBlock) error
	BanffProposalBlock(*BanffProposalBlock) error
	BanffStandardBlock(*BanffStandardBlock) error

	OdysseyAbortBlock(*OdysseyAbortBlock) error
	OdysseyCommitBlock(*OdysseyCommitBlock) error
	OdysseyProposalBlock(*OdysseyProposalBlock) error
	OdysseyStandardBlock(*OdysseyStandardBlock) error
	OdysseyAtomicBlock(*OdysseyAtomicBlock) error
}
