package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"image/color"
)

type Code struct {
	title string
	code  string
}

func (Code) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {

	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return Background{Color: color.NRGBA{R: 0xFA, G: 0xEA, B: 0xEF, A: 0xFF}}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(10).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Label(th, unit.Sp(18), "Title").Layout(gtx)
					})
				}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Bottom: unit.Dp(10),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.Label(th, unit.Sp(32), "123456").Layout(gtx)
						})
					})
				}))
			})
	})

}

type CodeView struct {
	list *layout.List
}

func NewCodeView() CodeView {
	list := layout.List{Axis: layout.Vertical, Alignment: layout.Middle}
	return CodeView{list: &list}
}

func (cv CodeView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	cv.list.Layout(gtx, 10, func(gtx layout.Context, index int) layout.Dimensions {
		return Code{}.Layout(gtx, th)
	})

	layout.SE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = 40
			return addIcon.Layout(gtx, color.NRGBA{R: 0xFF, A: 0xFF})
		})
	})

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}
