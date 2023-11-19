package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"image/color"
)

type IconButton struct{ size int }

func (i IconButton) Layout(gtx layout.Context, icon *widget.Icon, click *widget.Clickable) layout.Dimensions {
	btnColor := color.NRGBA{R: 0xDD, G: 0xDD, B: 0xDD, A: 0xFF}
	if click.Hovered() {
		btnColor = color.NRGBA{G: 0xFF, A: 0xFF}
	}
	if click.Pressed() {
		btnColor = color.NRGBA{R: 0xFF, A: 0xFF}
	}
	gtx.Constraints.Max.X = i.size
	gtx.Constraints.Min.X = i.size
	return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return icon.Layout(gtx, btnColor)
	})
}
