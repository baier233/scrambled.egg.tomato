package mylogger

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

var (
	red = color.New(color.FgRed).SprintFunc()
)

func Logf(format string, v ...any) {
	log.Printf(red(fmt.Sprintf(format, v)))
}

func Log(data string) {
	log.Println(red(data))
}
