package panels

import (
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/utils"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
)

func SettingsPanel(_ fyne.Window) fyne.CanvasObject {
	fix := widget.NewButton("修复", func() {
		err := os.Remove(utils.GetJreBinPath() + "\\winmm.dll")

		if err != nil {
			mylogger.LogErr("修复", err)
			return
		}
		mylogger.Log("已执行修复.")
	})

	return container.NewBorder(nil,
		container.NewVBox(fix, Line), nil, nil, nil)
}
