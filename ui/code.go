package ui

import (
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/dawenxi-tech/2fa/storage"
	"github.com/mazznoer/colorgrad"
	"image/color"
	"math"
	"slices"
	"time"
)

type Cell interface {
	Layout(gtx layout.Context, th *material.Theme) layout.Dimensions
}

type ToolCell struct {
	click widget.Clickable
	ctrl  *Controller

	text string
	icon *widget.Icon
}

func (cell *ToolCell) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if cell.click.Clicked() {
		if cell.text == "ADD CODE" {
			cell.ctrl.page = newAddView()
		} else {
			cell.ctrl.page = newSettingsView()
		}
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	var c = color.NRGBA{R: 0x81, G: 0x81, B: 0x81, A: 0xFF}
	dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return cell.click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return Background{Color: color.NRGBA{R: 0xFA, G: 0xEA, B: 0xEF, A: 0xFF}}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:    unit.Dp(40),
							Bottom: unit.Dp(40),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Alignment: layout.Middle}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return cell.icon.Layout(gtx, c)
								})
							}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								label := material.Label(th, unit.Sp(20), cell.text)
								label.Color = c
								return label.Layout(gtx)
							}))
						})
					})
				})
		})
	})
	return dims
}

type CodeCell struct {
	click  widget.Clickable
	id     string
	name   string
	secret string
	edit   bool
	input  *widget.Editor
	ctrl   *Controller

	delete widget.Clickable
}

func (c *CodeCell) initInput() {
	if c.edit && c.input == nil {
		c.input = &widget.Editor{SingleLine: true, Submit: true}
		c.input.SetText(c.name)
	}
}

func (c *CodeCell) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	c.processEvent(gtx)
	c.initInput()
	c.onSubmit(gtx)

	backgroundColor := color.NRGBA{R: 0xFA, G: 0xEA, B: 0xEF, A: 0xFF}
	layoutFn := ButtonLayoutStyle{CornerRadius: 4, Background: backgroundColor, Button: &c.click}.Layout
	if c.edit {
		layoutFn = Background{backgroundColor}.Layout
	}

	dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layoutFn(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(10).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					if c.edit {
						editor := material.Editor(th, c.input, "")
						editor.TextSize = unit.Sp(14)
						editor.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
						return editor.Layout(gtx)
					}
					label := material.Label(th, unit.Sp(14), c.name)
					label.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
					return label.Layout(gtx)
				})
			}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Bottom: unit.Dp(10),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						code := tryGetFA(c.secret)
						label := material.Label(th, unit.Sp(30), code)
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

func (c *CodeCell) processEvent(gtx layout.Context) {
	if c.delete.Clicked() {
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	if c.click.Clicked() {
		// copy code
		code := tryGetFA(c.secret)
		clipboard.WriteOp{Text: code}.Add(gtx.Ops)
	}
}

func (c *CodeCell) onSubmit(gtx layout.Context) {
	if c.input == nil {
		return
	}
	for _, event := range c.input.Events() {
		switch e := event.(type) {
		case widget.SubmitEvent:
			c.name = e.Text
			c.input = &widget.Editor{SingleLine: true, Submit: true}
			c.input.SetText(e.Text)
			return
		}
	}
}

type CodeView struct {
	list   layout.List
	edit   widget.Clickable
	ok     widget.Clickable
	cancel widget.Clickable
	add    widget.Clickable

	isEdit bool
	valid  bool

	cells    []Cell
	deleteId string
}

func newCodeView() *CodeView {
	list := layout.List{Axis: layout.Vertical, Alignment: layout.Middle}
	return &CodeView{list: list, edit: widget.Clickable{}}
}

func (cv *CodeView) Layout(gtx layout.Context, th *material.Theme, ctrl *Controller) layout.Dimensions {
	if !cv.valid {
		cv.reloadCodes(ctrl)
	} else if cv.deleteId != "" {
		cv.deleteCell(cv.deleteId)
		cv.deleteId = ""
	}

	if len(cv.cells) > 0 {
		op.InvalidateOp{At: time.Now().Add(time.Second * 5)}.Add(gtx.Ops)
	}

	if cv.add.Clicked() {
		ctrl.page = newAddView()
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	if cv.edit.Clicked() {
		cv.isEdit = !cv.isEdit
		if cv.isEdit {
			cv.cells = append(cv.cells, &ToolCell{ctrl: ctrl, text: "ADD CODE", icon: addIcon}, &ToolCell{ctrl: ctrl, text: "SETTINGS", icon: addIcon})
		} else {
			cv.cells = cv.cells[:len(cv.cells)-1]
		}
	}

	if cv.ok.Clicked() {
		cv.isEdit = false
		cv.syncCode()
		cv.valid = false
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	if cv.cancel.Clicked() {
		cv.isEdit = false
		cv.valid = false
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	if len(cv.cells) > 0 {
		cv.list.Layout(gtx, len(cv.cells), func(gtx layout.Context, index int) layout.Dimensions {
			if cell, ok := cv.cells[index].(*CodeCell); ok {
				cell.edit = cv.isEdit
			}
			return cv.cells[index].Layout(gtx, th)
		})
	}

	if cv.isEdit {
		layout.SE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Right:  unit.Dp(20),
					Bottom: unit.Dp(20),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return IconButton{size: 60}.Layout(gtx, okIcon, &cv.ok)
				})
			}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Right:  unit.Dp(20),
					Bottom: unit.Dp(20),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return IconButton{size: 60}.Layout(gtx, cancelIcon, &cv.cancel)
				})
			}))
		})
	} else if len(cv.cells) > 0 {
		layout.SE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return IconButton{size: 60}.Layout(gtx, editIcon, &cv.edit)
			})
		})
	} else {
		var c = color.NRGBA{R: 0x81, G: 0x81, B: 0x81, A: 0xFF}
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return cv.add.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    unit.Dp(40),
					Bottom: unit.Dp(40),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return addIcon.Layout(gtx, c)
						})
					}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Label(th, unit.Sp(20), "ADD CODE")
						label.Color = c
						return label.Layout(gtx)
					}))
				})
			})
		})
	}
	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

func (cv *CodeView) reloadCodes(ctrl *Controller) {
	codes := storage.LoadCodes()
	var cells []Cell
	for _, v := range codes {
		cells = append(cells, &CodeCell{id: v.ID, name: v.Name, secret: v.Secret.Val(), ctrl: ctrl})
	}
	if cv.isEdit {
		cells = append(cells, cv.cells[len(cv.cells)-2:]...)
	}
	cv.cells = cells
	cv.valid = true
}

func (cv *CodeView) deleteCell(id string) {
	cv.cells = slices.DeleteFunc(cv.cells, func(cell Cell) bool {
		if v, ok := cell.(*CodeCell); ok {
			return v.id == id
		}
		return false
	})
}

func (cv *CodeView) syncCode() {
	var codes []storage.Code
	for _, cell := range cv.cells {
		if v, ok := cell.(*CodeCell); ok {
			codes = append(codes, storage.Code{ID: v.id, Name: v.name})
		}
	}
	storage.SyncCode(codes)
}

var codeColorGradient = func() func() color.NRGBA {
	var gradient = func() colorgrad.Gradient {
		grad, _ := colorgrad.NewGradient().
			HtmlColors("#00FF00", "#FF0000").
			Build()
		return grad
	}()
	colors := gradient.ColorfulColors(6)
	return func() color.NRGBA {
		c := colors[time.Now().Second()%30/5]
		r, g, b := c.R*math.MaxUint8, c.G*math.MaxUint8, c.B*math.MaxUint8
		return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xFF}
	}
}()
