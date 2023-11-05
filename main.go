package main

import (
	"ScrambledEggwithTomato/data"
	"ScrambledEggwithTomato/localserver"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/panels"
	"ScrambledEggwithTomato/tm"
	"log"

	"github.com/fatih/color"
)

func preInit() {
	panels.Panels = map[string]panels.Panel{
		"注入": {"注入", panels.ModInjectPanel},
		"开端": {"开端", panels.ClientLaunchPanel},
	}
	panels.PanelIndex = map[string][]string{
		"": {"注入", "开端"}}

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

// go build -a -ldflags "-s -w"
func main() {

	// err := clientlauncher.InjectDllIntoMinecraft()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	red := color.New(color.FgHiRed, color.Bold).SprintFunc()
	log.SetPrefix("[" + red(" 炒.西红柿.鸡蛋 ") + "] ")
	preInit()
	data.Init()
	panels.Init()
	tm.IsDark = true
	panels.MyApp.Settings().SetTheme(&tm.MyTheme{})
	mylogger.Log("工具箱启动...")
	go localserver.BeginListen()
	defer mylogger.Log("工具箱关闭...")
	panels.Window.ShowAndRun()

}
