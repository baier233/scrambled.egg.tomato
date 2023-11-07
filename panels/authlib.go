package panels

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/utils"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Process struct {
	name string
	pid  uint32
}

type ProcessList struct {
	processes []Process
}

func getMinecrafts() (*ProcessList, error) {
	var processList ProcessList
	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, global.ErrorCreatCreateToolhelp32SnapshotFailed
	}

	defer syscall.CloseHandle(snapshot)

	var processEntry syscall.ProcessEntry32

	processEntry.Size = uint32(unsafe.Sizeof(processEntry))

	for syscall.Process32Next(snapshot, &processEntry) == nil {
		exeName := syscall.UTF16ToString(processEntry.ExeFile[:])
		if strings.Compare(exeName, "javaw.exe") != 0 || strings.Compare(exeName, "javaw.exe") != 0 {
			continue
		}

		cmdline, err := utils.GetCmdline(processEntry.ProcessID)
		if err != nil {
			continue
		}

		if strings.Contains(strings.ToUpper(cmdline), strings.ToUpper("-DlauncherControlPort")) {
			continue
		}

		processList.processes = append(processList.processes, Process{
			exeName,
			uint32(processEntry.ProcessID),
		})
		continue
	}
	return &processList, nil
}

var tempProcessList ProcessList

func AuthlibPanel(_ fyne.Window) fyne.CanvasObject {
	tempProcessList, err := getMinecrafts()
	if err != nil {
		mylogger.LogErr("获取minecraft进程", err)
	}

	selectedLabel := widget.NewLabel("No selection")
	data := binding.BindStringList(&tempModList)
	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	list.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < data.Length() {
			item, _ := data.GetItem(id)
			str, _ := item.(binding.String).Get()
			selectedLabel.SetText(str)
		} else {
			selectedLabel.SetText("No selection")
		}
	}
	for i, p := range tempProcessList.processes {
		data.Append(p.name + "  PID:" + strconv.Itoa(i))
	}
	inject := widget.NewButton("注入", func() {

		mylogger.Log("已注入 " + selectedLabel.Text)
	})

	return container.NewBorder(nil,
		container.NewVBox(inject, Line), nil, nil, list)
}
