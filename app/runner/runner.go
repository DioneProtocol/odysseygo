// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package runner

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/dioneprotocol/dionego/app"
	"github.com/dioneprotocol/dionego/app/process"
	"github.com/dioneprotocol/dionego/node"
)

// Run an DioneGo node.
func Run(nodeConfig node.Config) {
	nodeApp := process.NewApp(nodeConfig) // Create node wrapper
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println(process.Header)
	}

	exitCode := app.Run(nodeApp)
	os.Exit(exitCode)
}
