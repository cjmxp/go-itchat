// +build !windows

package utils

import "fmt"

func printWhite(s string) {
	fmt.Printf("\033[%d;%dm%s\033[0m", 37, 47, s)
}
