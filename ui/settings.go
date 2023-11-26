package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/dawenxi-tech/2fa/storage"
)

type Checkbox struct {
	selected bool
}

type SettingsView struct {
	showTray widget.Bool
	exit     widget.Bool
}

func newSettingsView() *SettingsView {
	conf := storage.LoadConfigure()
	sv := &SettingsView{}
	sv.showTray.Value = conf.ShowTray
	sv.exit.Value = conf.ExitWhenWindowClose
	return sv
}

func (s *SettingsView) Layout(gtx layout.Context, th *material.Theme, ctrl *Controller) layout.Dimensions {
	if s.exit.Changed() {
		s.saveConfigure(ctrl)
	}
	if s.showTray.Changed() {
		s.saveConfigure(ctrl)
	}
	layout.Inset{Top: 40, Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: 30}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(20).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: 30}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Right: 20}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Switch(th, &s.showTray, "show tray").Layout(gtx)
					})
				}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Label(th, unit.Sp(14), "Show Tray").Layout(gtx)
				}))
			})
		}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(20).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: 30}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Right: 20}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Switch(th, &s.exit, "Exit when window close").Layout(gtx)
					})
				}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Label(th, unit.Sp(14), "Exit when window close").Layout(gtx)
				}))
			})
		}))
	})

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

func (s *SettingsView) saveConfigure(ctrl *Controller) {
	conf := storage.LoadConfigure()
	conf.ExitWhenWindowClose = s.exit.Value
	conf.ShowTray = s.showTray.Value
	storage.SaveConfigure(conf)
	ctrl.win.configureChanged()
}
