package panels

import (
	"ScrambledEggwithTomato/clientlauncher"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var CheckGroup = widget.NewCheckGroup([]string{"开启开端", "开启proxy"}, onSelect)
var selectEntryStr = []string{"BaierCL"}

var EnabledProxy = false

func onSelect(s []string) {
	switch len(s) {
	case 0:
		{
			if clientlauncher.EnabledCL {
				clientlauncher.OnCloseCL()
			}
		}
	case 1:
		{
			if s[0] == "开启开端" {
				clientlauncher.OnEnableCL()
			} else {
				if clientlauncher.EnabledCL {
					clientlauncher.OnCloseCL()
				}
				EnabledProxy = true
			}
			break
		}
	case 2:
		{
			clientlauncher.OnEnableCL()
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
