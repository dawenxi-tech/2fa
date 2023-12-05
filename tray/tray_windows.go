//go:build windows
// +build windows

package tray

import (
	_ "embed"
	"fyne.io/systray"
	"github.com/dawenxi-tech/2fa/storage"
	"github.com/xlzd/gotp"
	"golang.design/x/clipboard"
	"sync"
	"time"
)

//go:embed 2fa-tray-win.ico
var iconData []byte

var once sync.Once

func show_tray() {
	once.Do(func() {
		go systray.Run(onReady, onExit)
	})
}

func dismiss_tray() {

}

func bring_window_to_front() {

}

func onExit() {

}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("2FA")

	var menus []*MenuItem

	codes := storage.LoadCodes()
	for _, code := range codes {
		menu := newCodeMenu(code.Name, code.Secret.Val())
		menu.show().listenClick()
		menus = append(menus, menu)
	}
	if len(codes) > 0 {
		systray.AddSeparator()
	}
	newMenuItem("Settings", EventShowSetting).show().listenClick()
	newMenuItem("Show Window", EventShowWindow).show().listenClick()
	newMenuItem("Quit", EventShowQuit).show().listenClick()

	refreshCodes(menus)
}

func change_application_activation_policy(i int) {}

const (
	MenuItemCode = iota + 1
	MenuItemMenu
)

type MenuItem struct {
	title  string
	secret string
	event  EventType
	typ    int
	*systray.MenuItem
	subMenu *systray.MenuItem
}

func (m *MenuItem) show() *MenuItem {
	if m.MenuItem == nil {
		switch m.typ {
		case MenuItemCode:
			menu := systray.AddMenuItem(m.title, "")
			totp := gotp.NewDefaultTOTP(m.secret)
			subMenu := menu.AddSubMenuItem(totp.Now(), "")
			m.MenuItem = menu
			m.subMenu = subMenu
		case MenuItemMenu:
			menu := systray.AddMenuItem(m.title, "")
			m.MenuItem = menu
		}
	}

	m.MenuItem.Show()
	return m
}

func (m *MenuItem) listenClick() {
	if m.MenuItem == nil {
		return
	}

	switch m.typ {
	case MenuItemCode:
		m.onCodeClick()
	case MenuItemMenu:
		m.onMenuClick()
	}
}

func (m *MenuItem) onCodeClick() {
	go func() {
		for {
			select {
			case <-m.MenuItem.ClickedCh:
			case <-m.subMenu.ClickedCh:
			}
			clipboard.Write(clipboard.FmtText, []byte(gotp.NewDefaultTOTP(m.secret).Now()))
		}
	}()
}

func (m *MenuItem) onMenuClick() {

	go func() {
		for {
			<-m.MenuItem.ClickedCh
			Event <- m.event
		}
	}()
}

func newCodeMenu(title string, secret string) *MenuItem {
	return &MenuItem{
		title:  title,
		secret: secret,
		typ:    MenuItemCode,
	}
}

func newMenuItem(title string, eventTyp EventType) *MenuItem {
	return &MenuItem{
		title: title,
		typ:   MenuItemMenu,
		event: eventTyp,
	}
}

func refreshCodes(menus []*MenuItem) {
	go func() {
		time.Sleep(time.Second * (time.Duration(30 - time.Now().Second()%30)))
		for _, m := range menus {
			if m.subMenu == nil {
				continue
			}
			code := gotp.NewDefaultTOTP(m.secret)
			m.subMenu.SetTitle(code.Now())
		}
	}()
}
