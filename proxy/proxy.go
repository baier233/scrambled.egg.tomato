package proxy

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var CurrentUUID uuid.UUID

func EstablishServer(data []string) error {
	defer func() {
		if r := recover(); r != nil {
			mylogger.Log("Recovered in EstablishServer", r)
		}
	}()
	if len(data) != 4 {
		return global.ErrorInternalIncorrectData
	}
	serverIp := data[0]
	serverPort := data[1]
	roleName := data[2]
	localPort := data[3]
	appName := syscall.StringToUTF16Ptr("Egg.Proxy.exe")

	commandLine := syscall.StringToUTF16Ptr(fmt.Sprintf("serverIp=%s serverPort=%s roleName=%s localPort=%s", serverIp, serverPort, roleName, localPort))

	var startupInfo windows.StartupInfo
	startupInfo.Cb = uint32(unsafe.Sizeof(startupInfo))

	var procInfo windows.ProcessInformation
	err := windows.CreateProcess(appName, commandLine, nil, nil, false, 0, nil, nil, &startupInfo, &procInfo)
	if err != nil {
		mylogger.LogErr("创建Proxy时", err)
	}
	global.CurrentServer = &procInfo

	return nil

}
