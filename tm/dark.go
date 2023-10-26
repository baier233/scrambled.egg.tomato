package tm

import (
	"ScrambledEggwithTomato/resources"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var IsDark bool

type MyDark struct{}

var _ fyne.Theme = (*MyDark)(nil)

func (m MyDark) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	//return theme.DefaultTheme().Color(n, v)
	if IsDark {
		return theme.DarkTheme().Color(n, v)
	}

	return theme.LightTheme().Color(n, v)
}
func (m MyDark) Icon(name fyne.ThemeIconName) fyne.Resource {
	if IsDark {
		return theme.DarkTheme().Icon(name)
	}

	return theme.LightTheme().Icon(name)
}

func (m MyDark) Font(style fyne.TextStyle) fyne.Resource {
	return &fyne.StaticResource{
		StaticName:    "HarmonyOS_Sans_SC_Regular.ttf",
		StaticContent: resources.FontHarmony,
	}
}

func (m MyDark) Size(name fyne.ThemeSizeName) float32 {
	if IsDark {
		return theme.DarkTheme().Size(name)
	}

	return theme.LightTheme().Size(name)
}
