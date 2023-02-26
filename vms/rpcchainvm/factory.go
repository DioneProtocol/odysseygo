// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rpcchainvm

import (
	"context"
	"fmt"
	"io"

	"github.com/dioneprotocol/dionego/snow"
	"github.com/dioneprotocol/dionego/utils/logging"
	"github.com/dioneprotocol/dionego/utils/resource"
	"github.com/dioneprotocol/dionego/vms"
	"github.com/dioneprotocol/dionego/vms/rpcchainvm/grpcutils"
	"github.com/dioneprotocol/dionego/vms/rpcchainvm/runtime"
	"github.com/dioneprotocol/dionego/vms/rpcchainvm/runtime/subprocess"

	vmpb "github.com/dioneprotocol/dionego/proto/pb/vm"
)

var _ vms.Factory = (*factory)(nil)

type factory struct {
	path           string
	processTracker resource.ProcessTracker
	runtimeTracker runtime.Tracker
}

func NewFactory(path string, processTracker resource.ProcessTracker, runtimeTracker runtime.Tracker) vms.Factory {
	return &factory{
		path:           path,
		processTracker: processTracker,
		runtimeTracker: runtimeTracker,
	}
}

func (f *factory) New(ctx *snow.Context) (interface{}, error) {
	config := &subprocess.Config{
		HandshakeTimeout: runtime.DefaultHandshakeTimeout,
	}

	// createStaticHandlers will send a nil ctx to disable logs
	// TODO: create a separate log file and no-op ctx
	if ctx != nil {
		config.Stderr = ctx.Log
		config.Stdout = ctx.Log
		config.Log = ctx.Log
	} else {
		config.Stderr = io.Discard
		config.Stdout = io.Discard
		config.Log = logging.NoLog{}
	}

	listener, err := grpcutils.NewListener()
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	status, stopper, err := subprocess.Bootstrap(
		context.TODO(),
		listener,
		subprocess.NewCmd(f.path),
		config,
	)
	if err != nil {
		return nil, err
	}

	clientConn, err := grpcutils.Dial(status.Addr)
	if err != nil {
		return nil, err
	}

	vm := NewClient(vmpb.NewVMClient(clientConn))
	vm.SetProcess(ctx, stopper, status.Pid, f.processTracker)

	f.runtimeTracker.TrackRuntime(stopper)

	return vm, nil
}
