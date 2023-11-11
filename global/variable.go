package global

import (
	"unsafe"

	"fyne.io/fyne/v2"
)

var Window fyne.Window
var CurrentServer unsafe.Pointer
var EnabledMod = false
var EnabledCL = false
var EnabledRemoveSrvMod = false

const (
	VersionNULL = iota
	Version4399
	Version163
	VersionMeatHook
)

var WPFVersion = VersionNULL

var CurrentUserInfo *UserInfo
