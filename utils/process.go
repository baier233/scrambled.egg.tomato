package utils

import (
	"reflect"
	"unsafe"

	"golang.org/x/sys/windows"
)

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
	err = windows.ReadProcessMemory(h, pbi.PebBaseAddress+0x20, // ntddk.h,ProcessParameters偏移32字节
		&d[0], uintptr(len(d)), nil)
	if err != nil {
		return "", err
	}

	var commandLine windows.NTUnicodeString
	Len := unsafe.Sizeof(commandLine)
	d = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&commandLine)),
		Len:  int(Len), Cap: int(Len)}))
	err = windows.ReadProcessMemory(h, uintptr(addr+0x70), // winternl.h,分析文件偏移
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
