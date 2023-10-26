package panels

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

const (
	emptyInputDataError  = "进来的数据不能是空的"
	nonexistentUserError = "用户是不存在的"
)

type Panel struct {
	Title string
	View  func(w fyne.Window) fyne.CanvasObject
}

var (
	Panels = map[string]Panel{
		"登录": {"登录", loginScreen},
		"注册": {"注册", registerScreen},
	}

	PanelIndex = map[string][]string{
		"": {"注册", "登录"},
	}
)

var (
	MyApp  = app.NewWithID("scrambled.egg.tomato")
	Window fyne.Window
	label  = widget.NewLabel("炒.西红柿.鸡蛋")
	Line   = canvas.NewLine(color.White)
)
