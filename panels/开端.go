package panels

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var CheckGroup = widget.NewCheckGroup([]string{"开启开端", "开启proxy"}, onSelect)
var selectEntryStr = []string{"BaierCL"}

var EnabledCL = false
var EnabledProxy = false

func onSelect(s []string) {
	switch len(s) {
	case 1:
		{
			if s[0] == "开启开端" {
				EnabledCL = true
			} else {
				EnabledProxy = true
			}
			break
		}
	case 2:
		{
			EnabledCL = true
			EnabledProxy = true
		}
	}
	if len(s) == 1 {
		if s[0] == "开启开端" {
			EnabledCL = true
		} else {
			EnabledProxy = true
		}
	}
}
func ClientLaunchPanel(_ fyne.Window) fyne.CanvasObject {

	selectEntry := widget.NewSelectEntry(selectEntryStr)
	selectItem := widget.NewForm(widget.NewFormItem("模式", selectEntry))

	selectEntry.SetText(selectEntryStr[0])

	CheckGroup.Horizontal = true
	return container.NewBorder(selectItem, CheckGroup, nil, nil, nil)
}
