package panels

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ClientLaunchPanel(_ fyne.Window) fyne.CanvasObject {
	button = widget.NewButton("开端!", func() {

	})
	return container.NewVBox(button)
}
