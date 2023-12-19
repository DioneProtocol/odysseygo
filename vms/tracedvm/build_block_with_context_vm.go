// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package tracedvm

import (
	"context"

	"go.opentelemetry.io/otel/attribute"

	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/DioneProtocol/odysseygo/snow/consensus/snowman"
	"github.com/DioneProtocol/odysseygo/snow/engine/snowman/block"
)

func (vm *blockVM) BuildBlockWithContext(ctx context.Context, blockCtx *block.Context) (snowman.Block, error) {
	if vm.buildBlockVM == nil {
		return vm.BuildBlock(ctx)
	}

	ctx, span := vm.tracer.Start(ctx, vm.buildBlockWithContextTag, oteltrace.WithAttributes(
		attribute.Int64("oChainHeight", int64(blockCtx.OChainHeight)),
	))
	defer span.End()

	return vm.buildBlockVM.BuildBlockWithContext(ctx, blockCtx)
}
