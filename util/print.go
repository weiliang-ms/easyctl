package util

import "fmt"

func PrintSuccessfulMsg(message string) {
	fmt.Printf("%c[1;40;32m[successful] %s%c[0m\n", 0x1B, message, 0x1B)
}

func PrintFailureMsg(message string) {
	fmt.Printf("%c[1;40;31m[failed] %s%c[0m\n", 0x1B, message, 0x1B)
}

func PrintTitleMsg(message string) {
	fmt.Printf("#### %s ####", message)
}

func PrintCloseServiceFailureMsg(message string) {
	fmt.Printf("%c[1;40;31m[failed] 关闭%s服务失败...%c[0m\n", 0x1B, message, 0x1B)
}

func PrintCloseServiceSuccessfulMsg(message string) {
	fmt.Printf("%c[1;40;32m[successful] 关闭%s服务成功...%c[0m\n", 0x1B, message, 0x1B)
}

func PrintRed(message string) string {
	return fmt.Sprintf("%c[1;40;31m[%s] %c[0m", 0x1B, message, 0x1B)
}

func PrintGreen(message string) string {
	return fmt.Sprintf("%c[1;40;32m[%s] %c[0m", 0x1B, message, 0x1B)
}

func PrintOrange(message string) string {
	return fmt.Sprintf("%c[1;40;33m[%s] %c[0m", 0x1B, message, 0x1B)
}

func PrintBlue(message string) string {
	return fmt.Sprintf("%c[1;40;34m[%s] %c[0m", 0x1B, message, 0x1B)
}

func PrintCyan(message string) string {
	return fmt.Sprintf("%c[1;40;36m[%s] %c[0m", 0x1B, message, 0x1B)
}
