package panels

import (
	"ScrambledEggwithTomato/clientlauncher"
	"ScrambledEggwithTomato/proxy"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var CheckGroup = widget.NewCheckGroup([]string{"开启开端", "开启proxy"}, onSelect)
var pCheckGroup unsafe.Pointer
var selectEntryStr = []string{"BaierCL"}

func ClientLaunchPanel(_ fyne.Window) fyne.CanvasObject {

	selectEntry := widget.NewSelectEntry(selectEntryStr)
	selectItem := widget.NewForm(widget.NewFormItem("模式", selectEntry))

	selectEntry.SetText(selectEntryStr[0])

	CheckGroup.Horizontal = true
	pCheckGroup = unsafe.Pointer(&CheckGroup)
	return container.NewBorder(selectItem, CheckGroup, nil, nil, nil)
}
func changeIt(strs []string) {

	check := (**widget.CheckGroup)(pCheckGroup)
	(*check).Selected = strs
}

func onSelect(s []string) {
	switch len(s) {
	case 0:
		{
			clientlauncher.OnCloseCL()
			proxy.OnCloseProxy()
		}
	case 1:
		{
			if s[0] == "开启开端" {
				s = clientlauncher.OnEnableCL(s)
				changeIt(s)
				proxy.OnCloseProxy()
			} else {
				clientlauncher.OnCloseCL()
				proxy.OnEnableProxy()
			}
			break
		}
	case 2:
		{
			s = clientlauncher.OnEnableCL(s)
			changeIt(s)
			proxy.OnEnableProxy()
		}
	}

}
