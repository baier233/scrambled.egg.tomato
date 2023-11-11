package tm

import (
	"ScrambledEggwithTomato/panels"
	"ScrambledEggwithTomato/resources"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var IsDark bool

type MyTheme struct{}

var _ fyne.Theme = (*MyTheme)(nil)

func (m MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {

	//return theme.DefaultTheme().Color(n, v)
	if IsDark {

		return theme.DarkTheme().Color(n, v)
	}

	return theme.LightTheme().Color(n, v)
}
func (m MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	if IsDark {

		return theme.DarkTheme().Icon(name)
	}

	return theme.LightTheme().Icon(name)
}

func (m MyTheme) Font(style fyne.TextStyle) fyne.Resource {

	if style.Bold {
		return &fyne.StaticResource{
			StaticName:    "HarmonyOS_Sans_SC_Bold.ttf",
			StaticContent: resources.HarmonyOS_Sans_SC_Bold,
		}

	} /*
		if style.Monospace {
			return &fyne.StaticResource{
				StaticName:    "HarmonyOS_Sans_SC_Regular.ttf",
				StaticContent: resources.HarmonyOS_Sans_SC_Regular,
			}
		}*/
	return &fyne.StaticResource{
		StaticName:    "HarmonyOS_Sans_SC_Regular.ttf",
		StaticContent: resources.HarmonyOS_Sans_SC_Regular,
	}
}

func (m MyTheme) Size(name fyne.ThemeSizeName) float32 {
	if IsDark {
		go func() {
			panels.Line.StrokeColor = color.White
		}()
		return theme.DarkTheme().Size(name)
	}
	go func() {
		panels.Line.StrokeColor = color.Gray16{0x9FFF}
	}()
	return theme.LightTheme().Size(name)
}
