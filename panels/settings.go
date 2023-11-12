package panels

import (
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/utils"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func SettingsPanel(_ fyne.Window) fyne.CanvasObject {

	ad := widget.NewRichTextWithText("购买幻影盾Java2Native虚拟化混淆器，保护你的应用安全 1720766910")
	ad2 := widget.NewRichTextWithText("幻影盾安全，有口碑，有保障，更新速度快 是您保护Java程序的不二之选")
	ad3 := widget.NewRichTextWithText("幻影盾有完善的License系统，让您的产品更安全地管理用户")

	fix := widget.NewButton("修复", func() {

		err := os.Remove(utils.GetJreBinPath() + "\\winmm.dll")

		if err != nil {
			mylogger.LogErr("修复", err)
			return
		}
		mylogger.Log("已执行修复.")
	})

	return container.NewBorder(nil,
		container.NewVBox(fix, Line), nil, nil, container.NewVBox(ad, ad2, ad3))
}
