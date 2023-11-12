package main

import "C"
import (
	"ScrambledEggwithTomato/data"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/panels"
	"ScrambledEggwithTomato/resources"
	"ScrambledEggwithTomato/tm"
	"github.com/fatih/color"
	"log"
)

func initLoginPanel() {
	panels.Panels =
		map[string]panels.Panel{
			"登录": {Title: "登录", View: panels.LoginScreen},
			"注册": {Title: "注册", View: panels.RegisterScreen},
		}
	panels.PanelIndex = map[string][]string{
		"": {"注册", "登录"},
	}
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
	initLoginPanel()
	panels.Init()
	tm.IsDark = true
	//fmt.Println(utils.GetHWID())

	panels.MyApp.Settings().SetTheme(&tm.MyTheme{})
	panels.Window.SetIcon(resources.IconResource)

	mylogger.Log("工具箱启动...")
	defer mylogger.Log("工具箱关闭...")

	panels.Window.ShowAndRun()
}
