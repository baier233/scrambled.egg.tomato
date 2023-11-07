package clientlauncher

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/proxy"
	"errors"

	"fyne.io/fyne/v2/dialog"
)

func OnEnableCL(group []string) []string {
	if !global.EnabledCL {
		if global.EnabledMod {
			dialog.ShowInformation("炒.西红柿.鸡蛋", "mod注入和开端只能开启一个，请关闭注入后再打开开端", global.Window)
			for i, s := range group {
				if s == "开启开端" {
					group = append(group[:i], group[i+1:]...)
				}
			}
			return group
		}
		global.EnabledCL = true
		mylogger.Log("已开启CL...")
		go ClientLaunchProcessor()
	}
	return group
}
func OnCloseCL() {
	global.EnabledCL = false
}

var lock = false

func ClientLaunchProcessor() {
	defer mylogger.Log("CL协程守护结束...")
	if lock {
		return
	}
	lock = true
	defer func() {
		lock = false
	}()
	mylogger.Log("CL协程守护中...")
	for global.EnabledCL {
		var serverData = NewServerData()
		err := InjectDllIntoMinecraft(serverData)
		if err == nil {
			mylogger.Log("Cl加载成功[1]...")

			if global.CurrentServer != nil {
				server := (*proxy.MinecraftProxyServer)(global.CurrentServer)
				server.CloseServer()
				global.CurrentServer = nil

				mylogger.Log("已强制结束proxyserver")
			}

			if proxy.EnabledProxy {

				data4proxy := make([]string, 4)
				data4proxy[0] = serverData.ServerIP
				data4proxy[1] = serverData.ServerPort
				data4proxy[2] = serverData.Username
				data4proxy[3] = "25565"
				go func() {
					err := proxy.EstablishServer(data4proxy)
					if err != nil {
						mylogger.Log("启动proxy时遇到不可预期的错误:" + err.Error())
					}
				}()

			}

		} else if !errors.Is(global.ErrorNonExistentMinecraftProcess, err) {
			mylogger.LogErr("注入CL", err)
		}
	}
}
