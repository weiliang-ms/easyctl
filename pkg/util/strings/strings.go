package strings

import (
	"fmt"
	"strings"
)

// ContainAll 字符串是否含所有字符
func ContainAll(str string, suffix ...string) bool {
	for _, v := range suffix {
		if strings.HasSuffix(str, v) {
			return true
		}
	}
	return false
}

func SplitIfContain(str string, contains []string) ([]string, error) {
	count := 0
	splitChar := ""
	for _, v := range contains {
		if strings.Contains(str, v) {
			count++
			splitChar = v
		}
	}

	switch count {
	case 1:
		return strings.Split(str, splitChar), nil
	default:
		return nil, fmt.Errorf("%s分割符不在%s内", str, contains)
	}
}
