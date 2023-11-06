package global

import "errors"

var (
	ErrorInjectFailed                        = errors.New("ErrorInjectFailed")
	ErrorNonExistentMinecraftProcess         = errors.New("ErrorNonExistentMinecraftProcess")
	ErrorCreatCreateToolhelp32SnapshotFailed = errors.New("ErrorCreatCreateToolhelp32SnapshotFailed")
	ErrorEmptyInputData                      = errors.New("ErrorEmptyInputData")
	ErrorInternalIncorrectData               = errors.New("ErrorInternalIncorrectData")
)
