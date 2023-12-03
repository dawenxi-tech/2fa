//go:build ((linux && !android) || freebsd) && !nowayland
// +build linux,!android freebsd
// +build !nowayland

package tray

import (
	_ "embed"
)

//go:embed 2fa-tray.ico
var iconData []byte

func show_tray() {

}

func dismiss_tray() {

}

func bring_window_to_front() {

}

func onExit() {

}

func onReady() {

}

func change_application_activation_policy(i int) {}
