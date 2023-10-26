package panels

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func registerScreen(_ fyne.Window) fyne.CanvasObject {
	usernameInput := widget.NewEntry()
	passwordInput := widget.NewPasswordEntry()
	activationCodeInput := widget.NewEntry()
	form := widget.NewForm(widget.NewFormItem("用户的名字", usernameInput), widget.NewFormItem("通过的单词", passwordInput), widget.NewFormItem("激活 代码", activationCodeInput))
	button = widget.NewButton("注册!", func() {

	})

	Line.StrokeWidth = 5
	return container.NewVBox(

		label,
		form,
		widget.NewSeparator(),
		button,
		Line,
	)
}
