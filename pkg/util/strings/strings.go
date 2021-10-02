package strings

import "strings"

// HasSuffix 字符串是否含有后缀
func HasSuffix(str string, suffix ...string) bool {
	for _, v := range suffix {
		if strings.HasSuffix(str, v) {
			return true
		}
	}
	return false
}
