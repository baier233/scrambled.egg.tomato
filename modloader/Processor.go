package modloader

import (
	"ScrambledEggwithTomato/mylogger"
)

var (
	EnablleMod = false
)

func OnEnableMod() {
	if !EnablleMod {

		EnablleMod = true
		mylogger.Log("已开启mod注入...")
	}
}

func OnCloseMod() {
	if EnablleMod {
		EnablleMod = false
		mylogger.Log("已关闭mod注入...")
	}
}
