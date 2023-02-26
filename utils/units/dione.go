// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package units

// Denominations of value
const (
	NanoDione  uint64 = 1
	MicroDione uint64 = 1000 * NanoDione
	Schmeckle uint64 = 49*MicroDione + 463*NanoDione
	MilliDione uint64 = 1000 * MicroDione
	Dione      uint64 = 1000 * MilliDione
	KiloDione  uint64 = 1000 * Dione
	MegaDione  uint64 = 1000 * KiloDione
)
