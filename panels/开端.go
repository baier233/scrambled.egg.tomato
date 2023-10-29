package panels

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var CheckGroup = widget.NewCheckGroup([]string{"开启开端", "开启proxy"}, func(s []string) { fmt.Println("selected", s) })
var selectEntryStr = []string{"BaierCL"}

func onSelect(s []string) {

}
func ClientLaunchPanel(_ fyne.Window) fyne.CanvasObject {

	selectEntry := widget.NewSelectEntry(selectEntryStr)
	selectItem := widget.NewForm(widget.NewFormItem("模式", selectEntry))
	selectEntry.SetText(selectEntryStr[0])

	CheckGroup.Horizontal = true
	return container.NewBorder(selectItem, CheckGroup, nil, nil, nil)
}
