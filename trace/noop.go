// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package trace

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/DioneProtocol/odysseygo/utils/constants"
)

var Noop Tracer = noOpTracer{
	Tracer: trace.NewNoopTracerProvider().Tracer(constants.AppName),
}

// noOpTracer is an implementation of trace.Tracer that does nothing.
type noOpTracer struct {
	trace.Tracer
}

func (n noOpTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return n.Tracer.Start(ctx, spanName, opts...)
}

func (noOpTracer) Close() error {
	return nil
}
