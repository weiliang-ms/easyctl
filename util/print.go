package util

import "fmt"

func PrintSuccessfulMsg(message string) {
	fmt.Printf("%c[1;40;32m[successful] %s%c[0m\n", 0x1B, message, 0x1B)
}

func PrintFailureMsg(message string) {
	fmt.Printf("%c[1;40;31m[failed] %s%c[0m\n", 0x1B, message, 0x1B)
}
