package ui

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/dawenxi-tech/2fa/storage"
	"github.com/mazznoer/colorgrad"
	"image/color"
	"math"
	"time"
)

type Cell interface {
	Layout(gtx layout.Context, th *material.Theme) layout.Dimensions
}

type AddCode struct {
	click widget.Clickable
	ctrl  *Controller
}

func (add *AddCode) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if add.click.Clicked() {
		add.ctrl.page = PageAdd
		op.InvalidateOp{}.Add(gtx.Ops)
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
							code := tryGetFA(c.code)
							label := material.Label(th, unit.Sp(32), code)
							label.Color = codeColorGradient()
							return label.Layout(gtx)
						})
					})
				}))
			})
	})

	if !c.edit {
		return dims
	}

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
	valid  bool

	cells []Cell
}

func newCodeView() CodeView {
	list := layout.List{Axis: layout.Vertical, Alignment: layout.Middle}
	return CodeView{list: list, edit: widget.Clickable{}}
}

func (cv *CodeView) Layout(gtx layout.Context, th *material.Theme, ctrl *Controller) layout.Dimensions {
	if !cv.valid {
		cv.reloadCodes()
	}

	btnColor := color.NRGBA{R: 0xDD, G: 0xDD, B: 0xDD, A: 0xFF}
	if cv.edit.Hovered() {
		btnColor = color.NRGBA{G: 0xFF, A: 0xFF}
	}
	if cv.edit.Pressed() {
		btnColor = color.NRGBA{R: 0xFF, A: 0xFF}
	}

	if len(cv.cells) > 0 {
		op.InvalidateOp{At: time.Now().Add(time.Second)}.Add(gtx.Ops)
	}

	if cv.edit.Clicked() {
		cv.isEdit = !cv.isEdit
		if cv.isEdit {
			cv.cells = append(cv.cells, &AddCode{ctrl: ctrl})
		} else {
			cv.cells = cv.cells[:len(cv.cells)-1]
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

func (cv *CodeView) reloadCodes() {
	codes := storage.LoadCodes()
	var cells []Cell
	for _, v := range codes {
		cells = append(cells, Code{title: v.Name, code: v.Secret.Val()})
	}
	if cv.isEdit {
		cells = append(cells, cv.cells[len(cv.cells)-1])
	}
	cv.cells = cells
	cv.valid = true
}

func _cond[T any](trueOrFalse bool, trueValue T, falseValue T) T {
	if trueOrFalse {
		return trueValue
	}
	return falseValue
}

var codeColorGradient = func() func() color.NRGBA {
	var gradient = func() colorgrad.Gradient {
		grad, _ := colorgrad.NewGradient().
			HtmlColors("#00FF00", "#FF0000").
			Build()
		return grad
	}()
	colors := gradient.ColorfulColors(30)
	return func() color.NRGBA {
		c := colors[time.Now().Second()%30]
		r, g, b := c.R*math.MaxUint8, c.G*math.MaxUint8, c.B*math.MaxUint8
		return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xFF}
	}
}()
