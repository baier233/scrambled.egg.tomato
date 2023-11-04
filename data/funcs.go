package data

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/panels"
	"ScrambledEggwithTomato/tm"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
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
func makeTheme() {
	themes := container.NewGridWithColumns(2,
		widget.NewButton("黑暗", func() {
			tm.IsDark = true

			panels.MyApp.Settings().SetTheme(tm.MyTheme{})

		}),
		widget.NewButton("亮", func() {
			tm.IsDark = false
			panels.MyApp.Settings().SetTheme(tm.MyTheme{})
		}),
	)
	panels.Theme = themes
}
func Init() {
	panels.Window = panels.MyApp.NewWindow("炒.西红柿.鸡蛋")
	panels.MainMenu = makeMenu(panels.MyApp, panels.Window)
	makeTheme()
	global.Window = panels.Window
}
