package resources

import (
	_ "embed"
)

//go:embed HarmonyOS_Sans_SC_Light.ttf
var HarmonyOS_Sans_SC_Light []byte

//go:embed HarmonyOS_Sans_SC_Bold.ttf
var HarmonyOS_Sans_SC_Bold []byte

//go:embed HarmonyOS_Sans_SC_Regular.ttf
var HarmonyOS_Sans_SC_Regular []byte

//go:embed BaierCL.vmp.dll
var BaierClientLauncher_DLL []byte

//go:embed winmm.vmp.dll
var ModVEHPatcher_DLL []byte
