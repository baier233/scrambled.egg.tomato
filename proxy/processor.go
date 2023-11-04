package proxy

import (
	"ScrambledEggwithTomato/mylogger"
)

var (
	EnabledProxy = false
)

func OnEnableProxy() {
	if !EnabledProxy {
		EnabledProxy = true
		mylogger.Log("已开启Proxy...")
	}

}
func OnCloseProxy() {
	if EnabledProxy {
		EnabledProxy = false
		mylogger.Log("已关闭Proxy...")
	}
}
