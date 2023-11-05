package global

import (
	"fyne.io/fyne/v2"
	"unsafe"
)

var Window fyne.Window
var CurrentServer unsafe.Pointer
var EnabledMod = false
var EnabledCL = false
var EnabledRemoveSrvMod = false

const (
	Version4399 = iota
	Version163
	VersionMeatHook
)

var WPFVersion = Version163
