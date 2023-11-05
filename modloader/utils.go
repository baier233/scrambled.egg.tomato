package modloader

import (
	"ScrambledEggwithTomato/mylogger"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

var is4399 bool = false

func getCustomModPath() string {
	absPath, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(err)
	}
	return absPath + "\\eggs"
}

func getMCDownloadPath() string {
	// 定义要读取的注册表项路径
	regPath := ""
	if is4399 {
		regPath = "Software\\Netease\\PC4399_MCLauncher\\DownloadPath"
	} else {
		regPath = "Software\\Netease\\MCLauncher\\DownloadPath"
	}
	// 构建命令行命令
	cmd := exec.Command("reg", "query", regPath)

	// 执行命令并获取输出
	output, err := cmd.Output()
	if err != nil {
		mylogger.Log(err.Error())
	}

	// 将输出转换为字符串并处理结果
	outputStr := string(output)
	outputLines := strings.Split(outputStr, "\n")
	for _, line := range outputLines {
		if strings.Contains(line, regPath) {
			// 找到包含注册表项路径的行，提取所需的值
			fields := strings.Fields(line)
			if len(fields) > 2 {
				value := fields[2]
				return value
			}
		}
	}

	return ""
}

func getJreBinPath() string {
	return getMCDownloadPath() + "\\ext\\jre-v64-220420\\jre8\\bin"
}

func getModsPath() string {
	return getMCDownloadPath() + "\\Game\\.minecraft\\mods"
}
