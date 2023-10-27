package main

import (
	"ScrambledEggwithTomato/data"
	"ScrambledEggwithTomato/panels"
	"ScrambledEggwithTomato/tm"
)

func preInit() {

	panels.Panels =
		map[string]panels.Panel{
			"登录": {"登录", panels.LoginScreen},
			"注册": {"注册", panels.RegisterScreen},
		}
	panels.PanelIndex = map[string][]string{
		"": {"注册", "登录"},
	}
}

func main() {
	preInit()
	data.Init()
	panels.Init()
	tm.IsDark = true
	panels.MyApp.Settings().SetTheme(&tm.MyTheme{})
	panels.Window.ShowAndRun()
}
