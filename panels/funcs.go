package panels

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const preferenceCurrentPanel = "登录"

var topWindow fyne.Window

func makeNav(setMenu func(panel Panel), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return PanelIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := PanelIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("啊？")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := Panels[uid]
			if !ok {
				//fyne.LogError("消失的界面: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			obj.(*widget.Label).TextStyle = fyne.TextStyle{}
		},
		OnSelected: func(uid string) {
			if t, ok := Panels[uid]; ok {
				a.Preferences().SetString(preferenceCurrentPanel, uid)
				setMenu(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentPanel, "登录")
		tree.Select(currentPref)
	}

	return container.NewBorder(nil, Theme, nil, nil, tree)
}

func Init() {
	Window.SetMainMenu(MainMenu)
	Window.SetMaster()
	topWindow = Window
	content := container.NewStack()
	setMenu := func(t Panel) {
		if fyne.CurrentDevice().IsMobile() {
			child := MyApp.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = Window
			})
			return
		}

		content.Objects = []fyne.CanvasObject{t.View(Window)}
		content.Refresh()
	}

	panel := container.NewBorder(
		nil, nil, nil, nil, content)

	split := container.NewHSplit(makeNav(setMenu, true), panel)
	split.Offset = 0
	Window.SetContent(split)

	Window.Resize(fyne.NewSize(820, 460))
}
