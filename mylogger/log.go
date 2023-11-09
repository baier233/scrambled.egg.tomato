package mylogger

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

var (
	red   = color.New(color.FgRed).SprintFunc()
	hired = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
)

func Logf(format string, v ...any) {
	log.Printf(red(fmt.Sprintf(format, v)))
}

func Log(data ...any) {
	log.Println(red(data))
}

func LogErr(prefix string, err error) {
	log.Println(green("在执行 ") + hired(prefix) + green(" 时出现预期之外的错误: ") + err.Error())
}
