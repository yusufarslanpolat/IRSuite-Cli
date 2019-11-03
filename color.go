package main

import (
	"fmt"
	ct "github.com/daviddengcn/go-colortext"
)

func Info(format string, args ...interface{}) {
	ct.Foreground(ct.Green, false)
	fmt.Printf(format, args...)
	ct.ResetColor()
}

func Warn(format string, args ...interface{}) {
	ct.Foreground(ct.Yellow, false)
	fmt.Printf(format, args...)
	ct.ResetColor()
}

func Error(format string, args ...interface{}) {
	ct.Foreground(ct.Red, false)
	fmt.Printf(format, args...)
	ct.ResetColor()
}
