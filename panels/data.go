package panels

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/login"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/utils"
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

var MainMenu *fyne.MainMenu
var Theme *fyne.Container

const preferenceCurrentPanel = "登录"

var topWindow fyne.Window

const (
	emptyInputDataError  = "进来的数据不能是空的"
	nonexistentUserError = "用户是不存在的"
)

type Panel struct {
	Title string
	View  func(w fyne.Window) fyne.CanvasObject
}

var Panels map[string]Panel
var PanelIndex map[string][]string

var (
	MyApp  = app.NewWithID("scrambled.egg.tomato")
	Window fyne.Window
	label  = widget.NewLabel("炒.西红柿.鸡蛋")
	Line   = canvas.NewLine(color.White)
)
var (
	currentUser = &CurrentUser{
		IsLoginIn: false,
		init:      false,
	}
)

type CurrentUser struct {
	init      bool
	user      *login.User
	IsLoginIn bool
}

func NewCurrentUser(theUser *login.User) *CurrentUser {
	return &CurrentUser{
		user:      theUser,
		IsLoginIn: theUser.Mark,
		init:      true,
	}
}
func initPanels(needMod bool) {

	Panels = map[string]Panel{
		"注入":      {"注入", ModInjectPanel},
		"开端":      {"开端", ClientLaunchPanel},
		"设置":      {"设置", SettingsPanel},
		"Authlib": {"Authlib", AuthlibPanel},
	}

	PanelIndex = map[string][]string{
		"": {"注入", "开端", "Authlib", "设置"}}

	slice := PanelIndex[""]
	if !needMod {
		slice = append(slice[:0], slice[1:]...)
	}
	PanelIndex[""] = slice
	return

}
func CheckWPFAndInitPanels() {
	wpfs := utils.GetWPFVersion()
	switch len(wpfs) {
	case 2:
		{
			dialog.ShowConfirm("检测到你同时安装了163和4399版本的网易盒子", "是否选择使用4399?", func(b bool) {
				if b {
					mylogger.Log("当前网易盒子版本4399")
					global.WPFVersion = global.Version4399
					return
				}
				mylogger.Log("当前网易盒子版本163")
				global.WPFVersion = global.Version163
			}, Window)
			break
		}
	case 1:
		{
			if wpfs[0] == "163" {
				global.WPFVersion = global.Version163
				break
			}
			global.WPFVersion = global.Version4399
			break
		}
	default:
		{
			mylogger.Log("你可能没安装网易盒子 mod注入功能将不可用! " + fmt.Sprint(wpfs))
			initPanels(false)
			return
		}
	}
	initPanels(true)
}
