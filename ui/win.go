package ui

import (
	"log"
	"os"
	"runtime"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/dawenxi-tech/2fa/storage"
	"github.com/dawenxi-tech/2fa/tray"
)

const _winWidth = 320
const _winHeight = 568

type Window struct {
	win  *app.Window
	ctrl *Controller
}

func NewWin() *Window {
	win := &Window{
		ctrl: newController(),
	}
	win.ctrl.win = win
	return win
}

func (w *Window) Run() {

	w.processTrayEvents()

	w.showWin()

	go func() {
		time.Sleep(time.Second * 2)
		w.showTray()
	}()

	app.Main()
}

func (w *Window) showTray() {
	conf := storage.LoadConfigure()
	if !conf.ShowTray {
		return
	}
	tray.ShowTray()
}

func (w *Window) loop() error {
	th := material.NewTheme()
	var ops op.Ops
	for {
		e := <-w.win.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			w.layout(gtx, th)
			e.Frame(gtx.Ops)
		}
	}
}

func (w *Window) layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return w.ctrl.Layout(gtx, th)
}

func (w *Window) showWin() {
	go func() {
		options := []app.Option{
			app.Title("2FA"),
			app.MinSize(unit.Dp(_winWidth), unit.Dp(_winHeight)),
			app.MaxSize(unit.Dp(_winWidth), unit.Dp(_winHeight)),
			app.Size(unit.Dp(_winWidth), unit.Dp(_winHeight)),
		}
		if runtime.GOOS == "darwin" {
			options = append(options, app.Decorated(false))
		}
		w.win = app.NewWindow(options...)
		w.win.Perform(system.ActionCenter)
		if err := w.loop(); err != nil {
			log.Fatal(err)
		}
		w.win = nil
	}()
}

func (w *Window) configureChanged() {
	conf := storage.LoadConfigure()
	if conf.ShowTray {
		w.showTray()
	} else {
		tray.DismissTray()
	}
}

func (w *Window) closeWin() {
	if w.win == nil {
		return
	}
	if _, ok := w.ctrl.page.(*SettingsView); ok {
		w.ctrl.page = newCodeView()
		return
	}
	conf := storage.LoadConfigure()
	if conf.ExitWhenWindowClose {
		os.Exit(1)
		return
	}
	w.win.Perform(system.ActionClose)
}

func (w *Window) processTrayEvents() {
	go func() {
		for {
			evt := <-tray.Event
			switch evt {
			case tray.EventShowSetting:
				w.ctrl.page = newSettingsView()
				if w.win != nil {
					w.win.Invalidate()
				}
				fallthrough
			case tray.EventShowWindow:
				if w.win == nil {
					w.showWin()
				}
				//tray.BringWindowToFront()
			default:
			}
		}
	}()
}

func (w *Window) exit() {
	os.Exit(0)
}
