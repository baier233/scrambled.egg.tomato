package main

import (
	"ScrambledEggwithTomato/clientlauncher"
	"ScrambledEggwithTomato/data"
	"ScrambledEggwithTomato/panels"
	"ScrambledEggwithTomato/tm"
	"fmt"
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

	err := clientlauncher.InjectDllIntoMinecraft()
	if err != nil {
		fmt.Println(err.Error())
	}
	preInit()
	data.Init()
	panels.Init()
	tm.IsDark = true
	panels.MyApp.Settings().SetTheme(&tm.MyTheme{})
	panels.Window.ShowAndRun()

}
