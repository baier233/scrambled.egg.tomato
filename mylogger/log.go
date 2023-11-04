package mylogger

import (
	"log"

	"github.com/fatih/color"
)

var (
	red = color.New(color.FgRed).SprintFunc()
)

func Log(data string) {
	log.Println(red(data))
}
