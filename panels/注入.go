package panels

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func ModInjectPanel(_ fyne.Window) fyne.CanvasObject {
	selectedLabel := widget.NewLabel("No selection")
	data := binding.BindStringList(&[]string{"Item 1", "Item 2", "Item 3"})

	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	add := widget.NewButton("Append", func() {
		val := fmt.Sprintf("Item %d", data.Length()+1)
		data.Append(val)
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

	button := widget.NewButton("注入！", func() {
		fmt.Println(selectedLabel.Text)
	})
	_ = button
	return container.NewBorder(nil, container.NewVBox(add, button), nil, nil, list)
}
