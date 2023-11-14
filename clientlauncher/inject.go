package clientlauncher

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/resources"
	"ScrambledEggwithTomato/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

//#include "inject.h"
//#include <stdlib.h>
//#cgo LDFLAGS: -L./ -ltest -lstdc++ -static
import "C"

var pidContainer = utils.NewStringContainer()

func pushData(cmdline string, serverData *ServerData) {
	var serverValue, portValue, usernameValue string
	serverPattern := "--server\\s+([^\\s]+)"
	serverRegexp := regexp.MustCompile(serverPattern)
	serverMatches := serverRegexp.FindStringSubmatch(cmdline)
	if len(serverMatches) > 1 {
		serverValue = serverMatches[1]
	}

	portPattern := "--port\\s+([^\\s]+)"
	portRegexp := regexp.MustCompile(portPattern)
	portMatches := portRegexp.FindStringSubmatch(cmdline)
	if len(portMatches) > 1 {
		portValue = portMatches[1]
	}

	usernamePattern := "--username\\s+([^\\s]+)"
	usernameRegexp := regexp.MustCompile(usernamePattern)
	usernameMatches := usernameRegexp.FindStringSubmatch(cmdline)
	if len(usernameMatches) > 1 {
		usernameValue = usernameMatches[1]
	}

	serverData.ServerIP = serverValue
	serverData.ServerPort = portValue
	serverData.Username = usernameValue
}

func InjectDllIntoMinecraft(serverData *ServerData) error {

	time.Sleep(time.Second / 2)

	var targetPid = uint32(0)
	{
		snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
		if err != nil {
			return global.ErrorCreatCreateToolhelp32SnapshotFailed
		}

		defer syscall.CloseHandle(snapshot)

		var processEntry syscall.ProcessEntry32

		processEntry.Size = uint32(unsafe.Sizeof(processEntry))

		for syscall.Process32Next(snapshot, &processEntry) == nil {

			if pidContainer.Contains(strconv.Itoa(int(processEntry.ProcessID))) {
				continue
			}

			exeName := syscall.UTF16ToString(processEntry.ExeFile[:])

			if strings.Compare(exeName, "javaw.exe") != 0 {
				continue
			}

			cmdline, err := utils.GetCmdline(processEntry.ProcessID)
			if err == nil {

				if strings.Contains(strings.ToUpper(cmdline), strings.ToUpper("-DlauncherControlPort")) {
					str := fmt.Sprintf("Process Name: %s, PID: %d ", exeName, processEntry.ProcessID)
					mylogger.Log("已找到Minecraft，" + str + " 准备执行注入操作...")
					pushData(cmdline, serverData)
					targetPid = processEntry.ProcessID
					pidContainer.Add(strconv.Itoa(int(targetPid)))
					break
				}
			}
		}
	}
	if targetPid == 0 {
		return global.ErrorNonExistentMinecraftProcess
	}

	result := C.inject(C.int(targetPid), (*C.char)(unsafe.Pointer(&resources.BaierClientLauncher_DLL[0])))
	if int(result) == 0 {
		return global.ErrorInjectFailed
	}
	return nil
	/*pidContainer.Add(strconv.Itoa(int(targetPid)))

		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		dPath := wd + "\\BaierCL.dll"

	here:
		_, err = os.Stat(dPath)
		if errors.Is(err, os.ErrNotExist) {
			file, err := os.Create(dPath)
			if err != nil {
				return err
			}
			_, err = file.Write(resources.BaierClientLauncher_DLL)
			if err != nil {
				return err
			}
			file.Close()
			goto here
		}

		kernel32 := windows.NewLazyDLL("kernel32.dll")
		pHandle, err := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, targetPid)
		if err != nil {
			return err
		}
		VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
		vAlloc, _, err := VirtualAllocEx.Call(uintptr(pHandle), 0, uintptr(len(dPath)+1), windows.MEM_RESERVE|windows.MEM_COMMIT, windows.PAGE_EXECUTE_READWRITE)

		bPtrDpath, err := windows.BytePtrFromString(dPath)
		if err != nil {
			return err
		}

		Zero := uintptr(0)
		err = windows.WriteProcessMemory(pHandle, vAlloc, bPtrDpath, uintptr(len(dPath)+1), &Zero)
		if err != nil {
			return err
		}

		LoadLibAddy, err := syscall.GetProcAddress(syscall.Handle(kernel32.Handle()), "LoadLibraryA")
		if err != nil {
			return err
		}

		tHandle, _, err := kernel32.NewProc("CreateRemoteThread").Call(uintptr(pHandle), 0, 0, LoadLibAddy, vAlloc, 0, 0)
		defer syscall.CloseHandle(syscall.Handle(tHandle))

		return nil*/
}
