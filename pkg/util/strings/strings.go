package strings

import "strings"

// ContainAll 字符串是否含所有字符
func ContainAll(str string, suffix ...string) bool {
	for _, v := range suffix {
		if strings.HasSuffix(str, v) {
			return true
		}
	}
	return false
}
