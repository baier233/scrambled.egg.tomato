package clientlauncher

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/resources"
	"fmt"
	"reflect"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

//#include "inject.h"
//#include <stdlib.h>
//#cgo LDFLAGS: -L./ -ltest -lstdc++ -lunwind -static
import "C"

var mapOfPID map[string]string

func InjectDllIntoMinecraft() error {

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
			exeName := syscall.UTF16ToString(processEntry.ExeFile[:])

			// if strings.Compare(exeName, "notepad.exe") == 0 {

			// 	fmt.Printf("Process Name: %s, PID: %d\n", exeName, processEntry.ProcessID)
			// 	targetPid = processEntry.ProcessID
			// 	break
			// }
			if strings.Compare(exeName, "javaw.exe") != 0 {
				continue
			}

			cmdline, err := GetCmdline(processEntry.ProcessID)
			if err == nil {
				fmt.Printf("Process Name: %s, PID: %d\n", exeName, processEntry.ProcessID)
				if strings.Contains(strings.ToUpper(cmdline), strings.ToUpper("-DlauncherControlPort")) {
					if mapOfPID[string(rune(targetPid))] == "Ok" {
						continue
					}
					mylogger.Log("已找到Minecraft，准备执行注入操作...")
					mapOfPID[string(rune(targetPid))] = "Ok"
					targetPid = processEntry.ProcessID
					break
				}
			}
		}
	}

	if targetPid == 0 {
		return global.ErrorNonExistentMinecraftProcess
	}
	fmt.Println(len(resources.BaierClientLauncher_DLL))
	result := C.inject(C.int(targetPid), (*C.char)(unsafe.Pointer(&resources.BaierClientLauncher_DLL[0])))
	if int(result) != 1 {
		return global.ErrortInjectFaield
	}
	return nil
	// 	fmt.Println(result)
	// 	return nil
	// 	wd, err := os.Getwd()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	dPath := wd + "\\BaierCL.dll"

	// here:
	// 	_, err = os.Stat(dPath)
	// 	if errors.Is(err, os.ErrNotExist) {
	// 		file, err := os.Create(dPath)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		_, err = file.Write(resources.BaierClientLauncher_DLL)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		file.Close()
	// 		goto here
	// 	}

	// 	kernel32 := windows.NewLazyDLL("kernel32.dll")
	// 	pHandle, err := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, targetPid)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	// 	vAlloc, _, err := VirtualAllocEx.Call(uintptr(pHandle), 0, uintptr(len(dPath)+1), windows.MEM_RESERVE|windows.MEM_COMMIT, windows.PAGE_EXECUTE_READWRITE)

	// 	bPtrDpath, err := windows.BytePtrFromString(dPath)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	Zero := uintptr(0)
	// 	err = windows.WriteProcessMemory(pHandle, vAlloc, bPtrDpath, uintptr(len(dPath)+1), &Zero)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	LoadLibAddy, err := syscall.GetProcAddress(syscall.Handle(kernel32.Handle()), "LoadLibraryA")
	// 	if err != nil {
	// 		return err
	// 	}

	// 	tHandle, _, err := kernel32.NewProc("CreateRemoteThread").Call(uintptr(pHandle), 0, 0, LoadLibAddy, vAlloc, 0, 0)
	// 	defer syscall.CloseHandle(syscall.Handle(tHandle))

	// 	return err
}

func GetCmdline(pid uint32) (string, error) {
	/* 翻译这个C++代码: https://stackoverflow.com/a/42341811/11844632 */
	if pid == 0 { // 系统进程,无法读取
		return "", nil
	}
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, pid)
	if err != nil {
		if e, ok := err.(windows.Errno); ok && e == windows.ERROR_ACCESS_DENIED {
			return "", nil // 没权限,忽略这个进程
		}
		return "", err
	}
	defer windows.CloseHandle(h)

	var pbi struct {
		ExitStatus                   uint32
		PebBaseAddress               uintptr
		AffinityMask                 uintptr
		BasePriority                 int32
		UniqueProcessId              uintptr
		InheritedFromUniqueProcessId uintptr
	}
	pbiLen := uint32(unsafe.Sizeof(pbi))
	err = windows.NtQueryInformationProcess(h, windows.ProcessBasicInformation, unsafe.Pointer(&pbi), pbiLen, &pbiLen)
	if err != nil {
		return "", err
	}

	var addr uint64
	d := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&addr)),
		Len:  8, Cap: 8}))
	err = windows.ReadProcessMemory(h, pbi.PebBaseAddress+32, // ntddk.h,ProcessParameters偏移32字节
		&d[0], uintptr(len(d)), nil)
	if err != nil {
		return "", err
	}

	var commandLine windows.NTUnicodeString
	Len := unsafe.Sizeof(commandLine)
	d = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&commandLine)),
		Len:  int(Len), Cap: int(Len)}))
	err = windows.ReadProcessMemory(h, uintptr(addr+112), // winternl.h,分析文件偏移
		&d[0], Len, nil)
	if err != nil {
		return "", err
	}

	cmdData := make([]uint16, commandLine.Length/2)
	d = *(*[]byte)(unsafe.Pointer(&cmdData))
	err = windows.ReadProcessMemory(h, uintptr(unsafe.Pointer(commandLine.Buffer)),
		&d[0], uintptr(commandLine.Length), nil)
	if err != nil {
		return "", err
	}
	return windows.UTF16ToString(cmdData), nil
}
