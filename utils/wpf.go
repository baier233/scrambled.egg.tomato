package utils

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"path/filepath"
)

func GetWPFVersion() []string {
	var data []string
	var value string
	reg4399 := `Software\Netease\PC4399_MCLauncher`
	reg163 := `Software\Netease\MCLauncher`
	key, err := registry.OpenKey(registry.CURRENT_USER, reg4399, registry.QUERY_VALUE)
	if err != nil {
		if registry.ErrNotExist.Is(err) {
			goto L163
		}
		mylogger.Log(err.Error())
	}
	defer key.Close()

	value, _, err = key.GetStringValue("DownloadPath") // 空字符串表示默认值
	if err != nil {
		if registry.ErrNotExist.Is(err) {
			goto L163
		}
	}
	if len(value) > 10 {
		data = append(data, "4399")
	}
L163:
	key, err = registry.OpenKey(registry.CURRENT_USER, reg163, registry.QUERY_VALUE)
	if err != nil {
		if registry.ErrNotExist.Is(err) {
			return data
		}
		mylogger.Log(err.Error())
	}
	defer key.Close()

	value, _, err = key.GetStringValue("DownloadPath") // 空字符串表示默认值
	if err != nil {
		if registry.ErrNotExist.Is(err) {
			return data
		}
	}
	if len(value) > 10 {
		data = append(data, "163")
	}
	return data
}

func GetCustomModPath() string {
	absPath, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(err)
	}
	return absPath + "\\eggs"
}

func GetMCDownloadPath() string {
	// 定义要读取的注册表项路径
	regPath := ""
	if global.WPFVersion == global.Version4399 {
		regPath = `Software\Netease\PC4399_MCLauncher`
	} else if global.WPFVersion == global.Version163 {
		regPath = `Software\Netease\MCLauncher`
	} else {
		mylogger.Log("你似乎没有安装网易盒子...")
		return ""
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.QUERY_VALUE)
	if err != nil {
		mylogger.Log(err.Error())
		return ""
	}
	defer key.Close()

	value, _, err := key.GetStringValue("DownloadPath")
	if err != nil {
		mylogger.Log(err.Error())
		return ""
	}

	return value
}

func GetJreBinPath() string {
	path := GetMCDownloadPath()
	if path != "" {
		return GetMCDownloadPath() + "\\ext\\jre-v64-220420\\jre8\\bin"
	}
	return ""
}

func GetModsPath() string {
	path := GetMCDownloadPath()
	if path != "" {

		return GetMCDownloadPath() + "\\Game\\.minecraft\\mods"
	}
	return ""
}
