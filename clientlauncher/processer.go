package clientlauncher

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"time"
)

var (
	EnabledCL = false
)

func OnEnableCL() {
	if !EnabledCL {
		EnabledCL = true
		mylogger.Log("已开启CL...")
		go ClientLaunchProcesser()
	}
	
}
func OnCloseCL() {
	EnabledCL = false
}
func ClientLaunchProcesser() {
	for EnabledCL {
		time.Sleep(time.Second)
		err := InjectDllIntoMinecraft()
		if err == nil {
			mylogger.Log("Cl加载成功...")
		}
		if err != global.ErrorNonExistentMinecraftProcess {
			mylogger.Log("出现不可预期的错误 : " + err.Error())
		}
	}
	mylogger.Log("已关闭CL...")
}
