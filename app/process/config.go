// (c) 2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package process

type Config struct {
	// If true, displays version and exits during startup
	DisplayVersionAndExit bool

	// String representation of app version
	VersionStr string

	// Path to the build directory
	BuildDir string

	// If true, run as a plugin
	PluginMode bool
}
