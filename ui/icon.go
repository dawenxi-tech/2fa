package ui

import (
	"gioui.org/widget"
	"github.com/dawenxi-tech/2fa/icon"
)

func loadIcon(data []byte) *widget.Icon {
	i, err := widget.NewIcon(data)
	if err != nil {
		panic(err)
	}
	return i
}

var closeIcon = loadIcon(icon.Close)

var circleIcon = loadIcon(icon.Circle)

var editIcon = loadIcon(icon.Edit)

var deleteIcon = loadIcon(icon.Delete)

var okIcon = loadIcon(icon.Ok)

var cancelIcon = loadIcon(icon.Cancel)

var addIcon = loadIcon(icon.Add)
