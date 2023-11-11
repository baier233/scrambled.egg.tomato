package resources

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed HarmonyOS_Sans_SC_Bold.ttf
var HarmonyOS_Sans_SC_Bold []byte

//go:embed HarmonyOS_Sans_SC_Regular.ttf
var HarmonyOS_Sans_SC_Regular []byte

//go:embed BaierCL.dll
var BaierClientLauncher_DLL []byte

//go:embed winmm.vmp.dll
var ModVEHPatcher_DLL []byte

//go:embed icon2.png
var icon []byte

var IconResource = &fyne.StaticResource{
	StaticName:    "icon.png",
	StaticContent: icon,
}
