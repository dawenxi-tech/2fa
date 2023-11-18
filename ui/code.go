package ui

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
)

type Cell interface {
	Layout(gtx layout.Context, th *material.Theme) layout.Dimensions
}

type AddCode struct {
	click widget.Clickable
}

func (add *AddCode) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if add.click.Clicked() {
		// goto add view
		fmt.Println("go to add view")
	}
	dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return add.click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return Background{Color: color.NRGBA{R: 0xFA, G: 0xEA, B: 0xEF, A: 0xFF}}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:    unit.Dp(40),
							Bottom: unit.Dp(40),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.Label(th, unit.Sp(32), "ADD CODE").Layout(gtx)
						})
					})
				})
		})
	})
	return dims
}

type Code struct {
	title string
	code  string
	edit  bool

	delete widget.Clickable
}

func (c Code) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {

	dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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

	ng := gtx
	ng.Constraints.Max = dims.Size
	layout.NE.Layout(ng, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = 32
			return c.delete.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return deleteIcon.Layout(gtx, color.NRGBA{R: 0xFF, A: 0xFF})
			})
		})
	})

	return dims
}

type CodeView struct {
	list layout.List
	edit widget.Clickable

	isEdit bool

	cells []Cell
}

func newCodeView() CodeView {
	list := layout.List{Axis: layout.Vertical, Alignment: layout.Middle}
	return CodeView{list: list, edit: widget.Clickable{}}
}

func (cv *CodeView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {

	btnColor := color.NRGBA{R: 0xDD, G: 0xDD, B: 0xDD, A: 0xFF}
	if cv.edit.Hovered() {
		btnColor = color.NRGBA{G: 0xFF, A: 0xFF}
	}
	if cv.edit.Pressed() {
		btnColor = color.NRGBA{R: 0xFF, A: 0xFF}
	}

	if cv.edit.Clicked() {
		cv.isEdit = !cv.isEdit
		if cv.isEdit {
			cv.cells = append(cv.cells, &AddCode{})
		}
	}

	cv.list.Layout(gtx, len(cv.cells), func(gtx layout.Context, index int) layout.Dimensions {
		return cv.cells[index].Layout(gtx, th)
	})

	layout.SE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = 60
			return cv.edit.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return editIcon.Layout(gtx, btnColor)
			})
		})
	})

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

func _cond[T any](trueOrFalse bool, trueValue T, falseValue T) T {
	if trueOrFalse {
		return trueValue
	}
	return falseValue
}
