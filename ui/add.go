package ui

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/dawenxi-tech/2fa/storage"
	"github.com/xlzd/gotp"
	"image"
	"image/color"
	"strings"
)

type AddView struct {
	editor    *widget.Editor
	codeInput *component.TextField
	applyBtn  *widget.Clickable
	cancelBtn *widget.Clickable
}

func newAddView() AddView {
	editor := &widget.Editor{
		SingleLine: true,
	}

	av := AddView{
		editor:    editor,
		applyBtn:  &widget.Clickable{},
		codeInput: &component.TextField{},
		cancelBtn: &widget.Clickable{},
	}

	return av
}

func (av AddView) Layout(gtx layout.Context, th *material.Theme, ctrl *Controller) layout.Dimensions {

	av.processEvents(gtx, ctrl)

	txt := av.codeInput.Text()

	code := tryGetFA(txt)

	layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return av.codeInput.Layout(gtx, th, "CODE OR URI")
			})
		}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Label(th, unit.Sp(30), code).Layout(gtx)
			})
		}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return material.Button(th, av.applyBtn, "ADD").Layout(gtx)
			})
		}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return material.Button(th, av.cancelBtn, "CANCEL").Layout(gtx)
			})
		}),
		)
	})

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

func (av AddView) processEvents(gtx layout.Context, ctrl *Controller) {
	if av.applyBtn.Clicked() {
		code := av.codeInput.Text()
		if secret := parseCodeOrUri(code); secret == "" {
			return
		}

		storage.InsertCode(code)

		ctrl.cv.valid = false
		ctrl.cv.isEdit = false
		ctrl.page = PageCode
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	if av.cancelBtn.Clicked() {
		ctrl.page = PageCode
		op.InvalidateOp{}.Add(gtx.Ops)
	}
}

func parseCodeOrUri(codeOrUri string) (secret string) {
	if strings.TrimSpace(codeOrUri) == "" {
		return ""
	}
	secret = codeOrUri
	parsed, _ := storage.ParseCode(codeOrUri)
	if parsed != nil {
		secret = parsed.Secret.Val()
	}
	defer func() {
		if x := recover(); x != nil {
			secret = ""
		}
	}()
	gotp.NewDefaultTOTP(secret).Now()
	return
}

func tryGetFA(code string) string {
	if secret := parseCodeOrUri(code); secret == "" {
		return "000000"
	} else {
		totp := gotp.NewDefaultTOTP(secret)
		return totp.Now()
	}
}

func drawBorder(ops *op.Ops, c color.NRGBA, width float32, x0, y0, x1, y1 int) {
	rrect := clip.RRect{Rect: image.Rectangle{
		Min: image.Pt(x0, y0),
		Max: image.Pt(x1, y1),
	}}
	paint.FillShape(ops, c,
		clip.Stroke{
			Path:  rrect.Path(ops),
			Width: width,
		}.Op(),
	)
}
