package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type AddView struct {
	editor *widget.Editor
}

func newAddView() AddView {
	editor := &widget.Editor{
		SingleLine: true,
	}

	av := AddView{
		editor: editor,
	}

	return av
}

func (av AddView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {

	layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.Editor(th, av.editor, "Please input the code.").Layout(gtx)
	})

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}
