package global

import (
	"fyne.io/fyne/v2"
	"golang.org/x/sys/windows"
)

var Window fyne.Window
var CurrentServer *windows.ProcessInformation
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

const ScrambledEggTomatoVersion = "1.0.3"
