package modloader

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
)

func OnEnableMod(binding *binding.Bool) {
	if !global.EnabledMod {
		if global.EnabledCL {
			dialog.ShowInformation("炒.西红柿.鸡蛋", "mod注入和开端只能开启一个，请关闭开端后再打开注入", global.Window)
			(*binding).Set(false)
			return
		}
		global.EnabledMod = true
		mylogger.Log("已开启mod注入...")
	}
}

func OnCloseMod() {
	if global.EnabledMod {
		global.EnabledMod = false
		mylogger.Log("已关闭mod注入...")
	}
}
