package main

import (
	"ScrambledEggwithTomato/panels"
	"ScrambledEggwithTomato/tm"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

const preferenceCurrentPanel = "currentPanel"

var topWindow fyne.Window

func main() {

	//panels.MyApp.Settings().SetTheme(theme.DarkTheme())
	tm.IsDark = true
	panels.MyApp.Settings().SetTheme(&tm.MyDark{})
	panels.Window = panels.MyApp.NewWindow("工具箱")
	panels.Window.SetMainMenu(makeMenu(panels.MyApp, panels.Window))
	panels.Window.SetMaster()
	topWindow = panels.Window
	content := container.NewStack()
	setMenu := func(t panels.Panel) {
		if fyne.CurrentDevice().IsMobile() {
			child := panels.MyApp.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = panels.Window
			})
			return
		}

		content.Objects = []fyne.CanvasObject{t.View(panels.Window)}
		content.Refresh()
	}

	panel := container.NewBorder(
		nil, nil, nil, nil, content)

	split := container.NewHSplit(makeNav(setMenu, true), panel)
	split.Offset = 0
	panels.Window.SetContent(split)

	panels.Window.Resize(fyne.NewSize(820, 460))
	panels.Window.ShowAndRun()
}

func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {

	cutShortcut := &fyne.ShortcutCut{Clipboard: w.Clipboard()}
	cutItem := fyne.NewMenuItem("✂", func() {
		shortcutFocused(cutShortcut, w)
	})
	cutItem.Shortcut = cutShortcut
	copyShortcut := &fyne.ShortcutCopy{Clipboard: w.Clipboard()}
	copyItem := fyne.NewMenuItem("复制", func() {
		shortcutFocused(copyShortcut, w)
	})
	copyItem.Shortcut = copyShortcut
	pasteShortcut := &fyne.ShortcutPaste{Clipboard: w.Clipboard()}
	pasteItem := fyne.NewMenuItem("粘贴", func() {
		shortcutFocused(pasteShortcut, w)
	})
	pasteItem.Shortcut = pasteShortcut

	main := fyne.NewMenu("主要的")
	mainMenu := fyne.NewMainMenu(
		main,
		fyne.NewMenu("编辑", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator()),
	)
	return mainMenu
}
func makeNav(setMenu func(panel panels.Panel), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return panels.PanelIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := panels.PanelIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("啊？")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := panels.Panels[uid]
			if !ok {
				fyne.LogError("消失的界面: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			obj.(*widget.Label).TextStyle = fyne.TextStyle{}
		},
		OnSelected: func(uid string) {
			if t, ok := panels.Panels[uid]; ok {
				a.Preferences().SetString(preferenceCurrentPanel, uid)
				setMenu(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentPanel, "登录")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("黑暗", func() {
			tm.IsDark = true
			panels.Line.StrokeColor = color.White
			panels.MyApp.Settings().SetTheme(tm.MyDark{})

		}),
		widget.NewButton("亮", func() {
			tm.IsDark = false
			panels.Line.StrokeColor = color.Gray16{0x9FFF}
			panels.MyApp.Settings().SetTheme(tm.MyLight{})
		}),
	)
	return container.NewBorder(nil, themes, nil, nil, tree)
}
