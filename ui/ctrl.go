package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
	"os"
)

type Page int

const (
	PageCode = iota
	PageAdd
	PageSettings
)

type Controller struct {
	win *Window

	page Page

	av *AddView
	cv *CodeView

	click widget.Clickable
}

func newController() *Controller {
	return &Controller{
		page: PageCode,
		av:   newAddView(),
		cv:   newCodeView(),
	}
}

func (ctrl *Controller) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {

	ctrl.processEvents(gtx)
	icon := circleIcon
	if ctrl.click.Hovered() {
		icon = closeIcon
	}

	switch ctrl.page {
	case PageCode:
		if ctrl.cv == nil {
			ctrl.cv = newCodeView()
		}
		ctrl.av = nil
		ctrl.cv.Layout(gtx, th, ctrl)
	case PageAdd:
		if ctrl.av == nil {
			ctrl.av = newAddView()
		}
		ctrl.cv = nil
		ctrl.av.Layout(gtx, th, ctrl)
	case PageSettings:

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
		os.Exit(0)
	}
}

func (ctrl *Controller) SwitchTo(page Page) {
	if ctrl.page == page {
		return
	}
	ctrl.page = page
}
