package clientlauncher

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/proxy"
)

var (
	EnabledCL = false
)

func OnEnableCL() {
	if !EnabledCL {
		EnabledCL = true
		mylogger.Log("已开启CL...")
		go ClientLaunchProcessor()
	}

}
func OnCloseCL() {
	EnabledCL = false
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
	for EnabledCL {
		var serverData = NewServerData()
		err := InjectDllIntoMinecraft(serverData)
		if err == nil {
			mylogger.Log("Cl加载成功...")

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

		} else if err != global.ErrorNonExistentMinecraftProcess {
			mylogger.Log("出现了不可预期的错误:" + err.Error())
		}
	}
}
