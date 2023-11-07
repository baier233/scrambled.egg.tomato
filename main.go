package main

import (
	"ScrambledEggwithTomato/data"
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/localserver"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/panels"
	"ScrambledEggwithTomato/resources"
	"ScrambledEggwithTomato/tm"
	"ScrambledEggwithTomato/utils"
	"fmt"
	"log"

	"fyne.io/fyne/v2/dialog"

	"github.com/fatih/color"
)

func initPanels(needMod bool) {

	panels.Panels = map[string]panels.Panel{
		"注入":      {"注入", panels.ModInjectPanel},
		"开端":      {"开端", panels.ClientLaunchPanel},
		"设置":      {"设置", panels.SettingsPanel},
		"Authlib": {"Authlib", panels.AuthlibPanel},
	}

	panels.PanelIndex = map[string][]string{
		"": {"注入", "开端", "Authlib", "设置"}}

	slice := panels.PanelIndex[""]
	if !needMod {
		slice = append(slice[:0], slice[1:]...)
	}
	panels.PanelIndex[""] = slice
	return

	panels.Panels =
		map[string]panels.Panel{
			"登录": {Title: "登录", View: panels.LoginScreen},
			"注册": {Title: "注册", View: panels.RegisterScreen},
		}
	panels.PanelIndex = map[string][]string{
		"": {"注册", "登录"},
	}
}
func checkWPFAndInitPanels() {
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
			}, panels.Window)
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

// go build -a -ldflags "-s -w"
func main() {
	// err := clientlauncher.InjectDllIntoMinecraft()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	red := color.New(color.FgHiRed, color.Bold).SprintFunc()
	log.SetPrefix("[" + red(" 炒.西红柿.鸡蛋 ") + "] ")

	data.Init()
	checkWPFAndInitPanels()
	panels.Init()
	tm.IsDark = true

	panels.MyApp.Settings().SetTheme(&tm.MyTheme{})
	panels.Window.SetIcon(resources.IconResource)

	mylogger.Log("工具箱启动...")
	go localserver.BeginListen()
	defer mylogger.Log("工具箱关闭...")

	panels.Window.ShowAndRun()
}
