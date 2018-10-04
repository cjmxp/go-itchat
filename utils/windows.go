// +build windows

package utils

import (
	"fmt"
	"syscall"
)

func printWhite(s string) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute")
	_, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(240)) //12 Red light
	fmt.Print(s)
	_, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7))
}
