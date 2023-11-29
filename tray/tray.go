package tray

type EventType int

const (
	EventShowSetting EventType = iota + 1
	EventShowWindow
)

const (
	ApplicationActivationPolicyAccessory int = 1
	ApplicationActivationPolicyRegular   int = 2
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

func ChangeApplicationActivationPolicy(i int) {
	change_application_activation_policy(i)
}
