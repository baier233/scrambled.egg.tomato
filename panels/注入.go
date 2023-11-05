package panels

import (
	"ScrambledEggwithTomato/modloader"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var tempLists []string
var EnabledInject = binding.NewBool()
var 非核心mod删除 = binding.NewBool()

func getCustomModPath() string {
	absPath, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(err)
	}
	return absPath + "\\eggs"
}

func ModInjectPanel(_ fyne.Window) fyne.CanvasObject {

	selectedLabel := widget.NewLabel("No selection")
	data := binding.BindStringList(&tempLists)

	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
	path := getCustomModPath()
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(getCustomModPath(), 0755)
		if err != nil {
			mylogger.Log("创建文件夹时出现预期之外的错误：" + err.Error())
		}
	}
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err = data.Append(info.Name())
			if err != nil {
				mylogger.Log("操作注入列表时出现预期之外的错误：" + err.Error())
				return err
			}
		}
		return nil
	})
	add := widget.NewButton("添加", func() {
		fileOpen := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
			if err != nil {
				mylogger.Log("选择文件时出现预期之外的错误：" + err.Error())
				return
			}
			if closer == nil {
				return
			}
			//fmt.Println(closer.URI().Path())

			err = utils.Copy(closer.URI().Path(), path+"\\"+closer.URI().Name())
			if err != nil {
				mylogger.Log("复制文件时出现预期之外的错误：" + err.Error())
			}

			strs, err := data.Get()
			if err != nil {
				return
			}

			for _, str := range strs {
				fmt.Println(str)
				if str == closer.URI().Name() {
					mylogger.Log("你添加了一个似乎已经添加过的mod名")
					return
				}
			}

			err = data.Append(closer.URI().Name())
			if err != nil {
				mylogger.Log("操作注入列表时出现预期之外的错误：" + err.Error())
				return
			}

		}, Window)
		fileOpen.SetFilter(storage.NewExtensionFileFilter([]string{".jar"}))
		fileOpen.Show()

	})

	del := widget.NewButton("删除", func() {

		if strings.EqualFold(selectedLabel.Text, "No selection") {
			return
		}
		lists, err := data.Get()
		if err != nil {
			return
		}
		var count = -1
		for i, s := range lists {
			if s == selectedLabel.Text {
				count = i
			}
		}
		if count == -1 {
			return
		}
		lists = append(lists[:count], lists[count+1:]...)
		err = data.Set(lists)
		if err != nil {
			return
		}
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
	EnabledInject.AddListener(binding.NewDataListener(func() {
		enabled, err := EnabledInject.Get()
		if err != nil {
			mylogger.Log("获取mod注入是否开启时出现了不可预期的错误：" + err.Error())
		}
		if enabled {
			modloader.OnEnableMod()
			return
		}
		modloader.OnCloseMod()
	}))
	check := widget.NewCheckWithData("开启注入", EnabledInject)
	check2 := widget.NewCheckWithData("删除非核心mod", 非核心mod删除)
	return container.NewBorder(nil,
		container.NewVBox(container.NewHBox(check, check2), container.NewVBox(add, Line), widget.NewSeparator(), container.NewVBox(del, Line)), nil, nil, list)
}
