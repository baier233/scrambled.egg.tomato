package clientlauncher

import (
	"ScrambledEggwithTomato/resources"
	"errors"
	"fmt"
	"github.com/zaneGittins/go-inject/inject"
	"golang.org/x/sys/windows"
	"os"
	"reflect"
	"strings"
	"syscall"
	"unsafe"
)

/*

 #cgo LDFLAGS: libReflectInject.a -lstdc++

#include <stdio.h>
#include <stdlib.h>
#include "library.h"

*/
import "C"

func InjectDllIntoMinecraft() error {

	snapshot := inject.CreateToolhelp32Snapshot(inject.TH32CS_SNAPPROCESS|inject.TH32CS_SNAPTHREAD, 0)
	var targetPid = uint32(0)
	var processEntry windows.ProcessEntry32
	processEntry.Size = uint32(unsafe.Sizeof(processEntry))
	_, err := inject.Process32Next(snapshot, &processEntry)
	if err != nil && err != windows.Errno(0) {
		fmt.Println(err)
		return err
	}

	for true {
		cmdline, err := GetCmdline(processEntry.ProcessID)
		if err == nil {
			if strings.Contains(strings.ToUpper(cmdline), strings.ToUpper("-DlauncherControlPort")) {
				fmt.Println(cmdline)
				fmt.Println()
				targetPid = processEntry.ProcessID
				break
			}
		}
		_, err = inject.Process32Next(snapshot, &processEntry)
		if err == windows.Errno(18) {
			break
		}
	}

	if targetPid == 0 {
		return errors.New("ErrorNotExistsMinecraftProcess")
	}
	reflectiveInject(int(targetPid), resources.BaierClientLauncher_DLL)
	return nil
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
	return err
}

//go:linkname injectFromBytes InjectFromBytes
//go:noescape
func injectFromBytes(pid C.int, pModuleBinary *C.char)

func reflectiveInject(pid int, moduleBinary []byte) {
	fmt.Println("Hi")
	res := unsafe.Pointer(&moduleBinary[0])
	result := (*C.char)(res)
	injectFromBytes(C.int(pid), result)
	fmt.Println("ok")
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
