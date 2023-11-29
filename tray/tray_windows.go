//go:build windows
// +build windows

package tray

import (
	_ "embed"
)

//go:embed 2fa-tray.png
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
