package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func CopyFile(src, dst string) error {
	sourceFilesStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFilesStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not correct", src)
	}

	source, err := os.Open(src)

	if err != nil {
		return err
	}

	defer source.Close()

	destination, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer destination.Close()
	_, err = io.Copy(destination, source)

	return err
}

func OpenFolderInExplorer(path string) {
	// 在资源管理器中打开文件夹
	cmd := exec.Command("explorer.exe", path)
	err := cmd.Start()
	if err != nil {
		return
	}
}
