package modloader

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/resources"
	"ScrambledEggwithTomato/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
		ReleaseRat()
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

func OnEnableRemoveSrvMod() {
	if !global.EnabledRemoveSrvMod {
		global.EnabledRemoveSrvMod = true
		mylogger.Log("已开启清理非核心mod...")
	}
}

func OnDisableRemoveSrvMod() {
	if global.EnabledRemoveSrvMod {
		global.EnabledRemoveSrvMod = false
		mylogger.Log("已关闭清理非核心mod...")
	}
}

func ReleaseRat() {
	file, err := os.Create(utils.GetJreBinPath() + "\\winmm.dll") // 创建或覆盖文件
	if err != nil {
		mylogger.LogErr("创建文件", err)
	}
	defer file.Close()

	_, err = file.Write(resources.ModVEHPatcher_DLL) // 写入数据
	if err != nil {
		mylogger.LogErr("写入文件", err)
	}
}
func InjectModProcessor() {
	mylogger.Log("注入mod中...")
	defer mylogger.Log("注入mod完成!")

	if global.EnabledRemoveSrvMod {
		ClearServerMods()
	}

	sourceDir := utils.GetCustomModPath()
	targetDir := utils.GetModsPath() // 修改为目标文件夹路径
	if sourceDir == "" || targetDir == "" {
		return
	}
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

func ClearServerMods() {
	root := utils.GetModsPath() // 替换为实际的文件夹路径
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(path, "@2@") {
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		mylogger.Logf("执行清理非核心mod出错：%v\n", err)
	} else {
		mylogger.Log("执行清理非核心mod完成.")
	}
}
