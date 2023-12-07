package main

import (
	"runtime/debug"

	"github.com/dawenxi-tech/2fa/ui"
)

// Version is initialized via ldflags or debug.BuildInfo.
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
