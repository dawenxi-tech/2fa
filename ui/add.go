package ui

import (
	"bytes"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/dawenxi-tech/2fa/storage"
	"github.com/dim13/otpauth/migration"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/xlzd/gotp"
	"golang.design/x/clipboard"
	"image"
	"image/color"
	"log/slog"
	"strings"
)

type AddView struct {
	codeInput *component.TextField
	applyBtn  *widget.Clickable
	cancelBtn *widget.Clickable

	codes  []string
	imag   image.Image
	isRead bool
}

func newAddView() *AddView {
	av := &AddView{
		applyBtn:  &widget.Clickable{},
		codeInput: &component.TextField{CharLimit: 2048},
		cancelBtn: &widget.Clickable{},
	}
	av.codeInput.Editor.SingleLine = true
	return av
}

func (av *AddView) Layout(gtx layout.Context, th *material.Theme, ctrl *Controller) layout.Dimensions {
	if !av.isRead {
		av.tryReadClipboard()
		av.isRead = true
	}
	av.processEvents(gtx, ctrl)
	if len(av.codes) > 0 && av.imag != nil {
		return av.layoutQR(gtx, th, ctrl)
	} else {
		return av.layoutTextField(gtx, th, ctrl)
	}
}

func (av *AddView) layoutQR(gtx layout.Context, th *material.Theme, ctrl *Controller) layout.Dimensions {

	var flexes []layout.FlexChild
	flexes = append(flexes, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: 20}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = 400
				gtx.Constraints.Max.X = 400
				return widget.Image{
					Src: paint.NewImageOp(av.imag),
					Fit: widget.Contain,
				}.Layout(gtx)
			})
		})
	}))

	for _, c := range av.codes {
		code := tryGetFA(c)
		flexes = append(flexes, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.Label(th, 20, code).Layout(gtx)
				})
			})
		}))
	}
	flexes = append(flexes, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return material.Button(th, av.applyBtn, "ADD").Layout(gtx)
		})
	}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return material.Button(th, av.cancelBtn, "CANCEL").Layout(gtx)
		})
	}))

	layout.Flex{Axis: layout.Vertical}.Layout(gtx, flexes...)

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

func (av *AddView) layoutTextField(gtx layout.Context, th *material.Theme, ctrl *Controller) layout.Dimensions {
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

func (av *AddView) tryReadClipboard() {
	// read clipboard image format
	data := clipboard.Read(clipboard.FmtImage)
	if len(data) == 0 {
		return
	}
	img, _, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		slog.With(slog.Any("err", err)).Error("error to decode image")
		return
	}
	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		slog.With(slog.Any("err", err)).Error("error to prepare zxing")
		return
	}
	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		slog.With(slog.Any("err", err)).Error("error to decode image")
		return
	}

	av.tryParseCode(result.GetText())
	if len(av.codes) > 0 {
		av.imag = img
	}
}

func (av *AddView) tryParseCode(txt string) {
	// try get 2fa
	if isCodeUriValid(txt) {
		av.codes = []string{txt}
		return
	}
	// try decode google uri
	av.tryParseGoogleAuthenticatorAppExportQr(txt)
}

func (av *AddView) tryParseGoogleAuthenticatorAppExportQr(qr string) {
	data, err := migration.Data(qr)
	if err != nil {
		return
	}
	param, err := migration.Unmarshal(data)
	if err != nil {
		return
	}
	var codes []string
	for _, p := range param.OtpParameters {
		uri := p.URL()
		if uri != nil {
			codes = append(codes, uri.String())
		}
	}
	av.codes = codes
}

func (av *AddView) processEvents(gtx layout.Context, ctrl *Controller) {
	if av.applyBtn.Clicked() {
		if len(av.codes) > 0 {
			for _, code := range av.codes {
				if !isCodeUriValid(code) {
					continue
				}
				storage.InsertCode(code)
			}
		} else {
			code := av.codeInput.Text()
			if !isCodeUriValid(code) {
				return
			}
			storage.InsertCode(code)
		}
		ctrl.page = PageCode
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	if av.cancelBtn.Clicked() {
		if len(av.codes) > 0 {
			av.codes = nil
			op.InvalidateOp{}.Add(gtx.Ops)
			return
		}
		ctrl.page = PageCode
		op.InvalidateOp{}.Add(gtx.Ops)
	}
}

func isCodeUriValid(uri string) bool {
	secret := parseCodeOrUri(uri)
	return secret != ""
}

func parseCodeOrUri(codeOrUri string) (secret string) {
	codeOrUri = strings.TrimSpace(codeOrUri)
	if codeOrUri == "" {
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
