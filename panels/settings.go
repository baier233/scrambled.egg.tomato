package panels

import (
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
)

func SettingsPanel(_ fyne.Window) fyne.CanvasObject {
	fix := widget.NewButton("修复", func() {
		err := os.Remove(utils.GetJreBinPath() + "\\winmm.dll")

		if err != nil {
			mylogger.Log("在修复时发生超出预期的错误:" + err.Error())
		}
		mylogger.Log("已执行修复.")
	})

	return container.NewBorder(container.NewVBox(fix, Line),
		nil, nil, nil, nil)
}
