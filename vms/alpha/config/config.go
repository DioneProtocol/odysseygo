// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package config

import "time"

// Struct collecting all the foundational parameters of the ALPHA
type Config struct {
	// Fee that is burned by every non-asset creating transaction
	TxFee uint64

	// Fee that is burned by every non-asset creating transaction
	ApricotPhase7TxFee uint64

	// Fee that must be burned by every asset creating transaction
	CreateAssetTxFee uint64

	// Time of the AP7 network upgrade
	ApricotPhase7Time time.Time
}

func (c *Config) IsApricotPhase7Activated(timestamp time.Time) bool {
	return !timestamp.Before(c.ApricotPhase7Time)
}

func (c *Config) GetTxFee(timestamp time.Time) uint64 {
	if c.IsApricotPhase7Activated(timestamp) {
		return c.ApricotPhase7TxFee
	}
	return c.TxFee
}
