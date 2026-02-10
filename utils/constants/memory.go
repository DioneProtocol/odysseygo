// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package constants

// PointerOverhead is used to approximate the memory footprint from allocating a
// pointer.
const PointerOverhead = 8

// OrionSampleSizeRatio defines the ratio of orions that will be present in the sample size list.
// The orion sample size is calculated as: floor(sampleSize * OrionSampleSizeRatio).
// For example, if sample size is 20 and ratio is 0.2, then orion sample size will be 4.
const OrionSampleSizeRatio float64 = 0
