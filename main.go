package main

import (
	"runtime/debug"

	"github.com/dawenxi-tech/2fa/ui"
)

// Version is initialized via ldflags or debug.BuildInfo.
// See: https://github.com/vorlif/xspreak/commit/8ff2092f126bf096a56062ac13126af8ece4eb71#diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34
var Version = ""

func init() {
	initVersionNumber()
}

func initVersionNumber() {
	// If already set via ldflags, the value is retained.
	if Version != "" {
		return
	}

	version := ""

	info, available := debug.ReadBuildInfo()

	if available {
		version = info.Main.Version
	} else {
		version = "dev"
	}

	Version = version
}

func main() {
	win := ui.NewWin()
	win.Run()
}
