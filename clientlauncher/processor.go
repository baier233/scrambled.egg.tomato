package clientlauncher

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/proxy"
	"ScrambledEggwithTomato/utils"
	"errors"
	"fyne.io/fyne/v2/dialog"
	"os"
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
	if _, err := os.Stat(utils.GetJreBinPath() + "\\winmm.dll"); err == nil {
		err := os.Remove(utils.GetJreBinPath() + "\\winmm.dll")
		if err != nil {
			mylogger.LogErr("开端-删除注入文件", err)
		}
	}

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
				ServerDataChan <- serverData
			}

		} else if !errors.Is(global.ErrorNonExistentMinecraftProcess, err) {
			mylogger.LogErr("注入CL", err)
		}
	}
}
