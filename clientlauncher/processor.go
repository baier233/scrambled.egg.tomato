package clientlauncher

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/proxy"
	"ScrambledEggwithTomato/utils"
	"errors"
	"fyne.io/fyne/v2/dialog"
	"golang.org/x/sys/windows"
	"os"
	"time"
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
		time.Sleep(time.Second / 2)
		var serverData = NewServerData()
		err := InjectDllIntoMinecraft(serverData)
		if err == nil {
			mylogger.Log("Cl加载成功[1]...")

			if global.CurrentServer != nil {
				err := windows.TerminateProcess(global.CurrentServer.Process, 1)
				if err != nil {
					mylogger.LogErr("结束proxyserver时", err)
				}
				windows.CloseHandle(global.CurrentServer.Process)
				windows.CloseHandle(global.CurrentServer.Thread)
				global.CurrentServer = nil
				mylogger.Log("已强制结束proxyserver")
			}

			if proxy.EnabledProxy {
				select {
				case ServerDataChan <- serverData:
				case <-time.After(time.Minute * 2):
					mylogger.LogErr("发送服务器数据时", errors.New("超时"))
				}
			}

		} else if !errors.Is(global.ErrorNonExistentMinecraftProcess, err) {
			mylogger.LogErr("注入CL", err)
		}
	}
}
