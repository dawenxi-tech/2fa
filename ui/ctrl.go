package ui

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
)

type Page interface {
	Layout(gtx layout.Context, th *material.Theme, ctrl *Controller) layout.Dimensions
}

type Controller struct {
	win   *Window
	page  Page
	click widget.Clickable
}

func newController() *Controller {
	return &Controller{
		page: newCodeView(),
	}
}

func (ctrl *Controller) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {

	ctrl.processEvents(gtx)
	icon := circleIcon
	if ctrl.click.Hovered() {
		icon = closeIcon
	}

	if ctrl.page != nil {
		ctrl.page.Layout(gtx, th, ctrl)
	}

	// close button
	layout.NW.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = 28
			return ctrl.click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return icon.Layout(gtx, color.NRGBA{R: 0xFC, G: 0x60, B: 0x5C, A: 0xFF})
			})
		})
	})

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

func (ctrl *Controller) processEvents(gtx layout.Context) {
	if ctrl.click.Clicked() {
		ctrl.win.closeWin()
		op.InvalidateOp{}.Add(gtx.Ops)
	}
}

func (ctrl *Controller) SwitchTo(page Page) {
	if ctrl.page == page {
		return
	}
	ctrl.page = page
}
