package ui

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"image/color"
)

type Background struct {
	Color color.NRGBA
}

func (b Background) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()
	paint.FillShape(gtx.Ops, b.Color, clip.Rect{Max: dims.Size}.Op())
	call.Add(gtx.Ops)
	return dims
}

type TitleBar struct{}

func (TitleBar) Layout(gtx layout.Context) layout.Dimensions {

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{
				Size: gtx.Constraints.Max,
			}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = 60
			gtx.Constraints.Max.Y = 60
			gtx.Constraints.Min.X = 60
			gtx.Constraints.Min.Y = 60
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = 32
				gtx.Constraints.Max.Y = 32
				return closeIcon.Layout(gtx, color.NRGBA{R: 0xFF, A: 0xFF})
			})
		}),
	)
}
