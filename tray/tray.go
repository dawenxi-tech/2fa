package tray

import (
	_ "embed"
)

//go:embed 2fa-tray.png
var iconData []byte

type EventType int

const (
	EventShowSetting EventType = iota + 1
	EventShowWindow
)

var Event = make(chan EventType, 2)

func ShowTray() {
	show_tray()
}

func DismissTray() {
	dismiss_tray()
}

func sendEvent(typ EventType) {
	Event <- typ
}

func BringWindowToFront() {
	bring_window_to_front()
}
