package modloader

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
)

func OnEnableMod(binding *binding.Bool) {
	if !global.EnabledMod {
		if global.EnabledCL {
			dialog.ShowInformation("炒.西红柿.鸡蛋", "mod注入和开端只能开启一个，请关闭开端后再打开注入", global.Window)
			(*binding).Set(false)
			return
		}
		global.EnabledMod = true
		mylogger.Log("已开启mod注入...")
	}
}

func OnCloseMod() {
	if global.EnabledMod {
		global.EnabledMod = false
		mylogger.Log("已关闭mod注入...")
	}
}

func InjectModProcessor() {
	sourceDir := getCustomModPath()
	targetDir := getModsPath() // 修改为目标文件夹路径
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".jar" {
			// 创建目标文件夹路径，如果还不存在的话
			targetPath := filepath.Join(targetDir, filepath.Base(path))
			if _, err := os.Stat(targetPath); err != nil {
				if os.IsNotExist(err) {
					os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
				} else {
					return err
				}
			}
			// 打开源文件和目标文件
			sourceFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer sourceFile.Close()
			targetFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer targetFile.Close()
			// 复制文件内容到目标文件
			_, err = io.Copy(targetFile, sourceFile)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking through files: %v\n", err)
		return
	}
}
