package panels

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"strings"
)

var tempLists []string

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

	add := widget.NewButton("添加", func() {
		fileopen := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println(err.Error())

				return
			}
			if closer == nil {
				return
			}
			fmt.Println(closer.URI().Path())
			err = data.Append(closer.URI().Name())
			if err != nil {
				fmt.Println(err.Error())
				return
			}

		}, Window)
		fileopen.SetFilter(storage.NewExtensionFileFilter([]string{".jar"}))
		fileopen.Show()

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

	return container.NewBorder(nil, container.NewVBox(container.NewVBox(add, Line), widget.NewSeparator(), container.NewVBox(del, Line)), nil, nil, list)
}
