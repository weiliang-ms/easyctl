package strings

import (
	"fmt"
	"strings"
)

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

// SubSlash 分割slash
func SubSlash(str string) []string {

	var splitSlice []string
	if !strings.Contains(str, "\\") && !strings.Contains(str, "/"){
		return []string{str}
	}

	if strings.Contains(str, "/") {
		for _, v := range strings.Split(str, "/") {
			splitSlice = append(splitSlice, v)
		}
	}

	ifContinue := false

	for _ , v := range splitSlice {
		if strings.Contains(v, "\\") {
			ifContinue=true
		}
	}

	if !ifContinue {
		return splitSlice
	}

	for _, v := range strings.Split(str, "\\") {
		splitSlice = append(splitSlice, v)
	}

	return splitSlice
}

// TrimPrefixAndSuffix = TrimPrefix & TrimSuffix
func TrimPrefixAndSuffix(str, fix string) string {
	if strings.HasPrefix(str, fix) && strings.HasSuffix(str, fix) {
		return strings.TrimSuffix(strings.TrimPrefix(str, fix), fix)
	}

	return str
}

// SubFileName 截取文件名称
func SubFileName(s string) string {

	nameSplitSlice := SubSlash(s)
	if len(nameSplitSlice) > 1 {
		return nameSplitSlice[len(nameSplitSlice)-1]
	}
	return s
}
