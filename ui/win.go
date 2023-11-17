package ui

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"image/color"
	"log"
	"os"
)

const _winWidth = 320
const _winHeight = 568

type Window struct {
	win *app.Window

	code CodeView
	add  AddView
}

func NewWin() *Window {
	return &Window{
		code: NewCodeView(),
		add:  newAddView(),
	}
}

func (w *Window) Run() {
	go func() {
		w.win = app.NewWindow(
			app.Decorated(false),
			app.Title("2FA"),
			app.MinSize(unit.Dp(_winWidth), unit.Dp(_winHeight)),
			app.MaxSize(unit.Dp(_winWidth), unit.Dp(_winHeight)),
			app.Size(unit.Dp(_winWidth), unit.Dp(_winHeight)),
		)
		w.win.Perform(system.ActionCenter)
		if err := w.loop(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
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

	//return w.code.Layout(gtx, th)

	w.add.Layout(gtx, th)

	layout.NW.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = 28
			return closeIcon.Layout(gtx, color.NRGBA{R: 0xFF, A: 0xFF})
		})
	})
	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}
