package util

import (
	"fmt"
	"runtime"
	"strings"
)

// AppendStringFromSlice 拼接字符串
func AppendStringFromSlice(slice []string, appendStr string) (s string) {

	for _, v := range slice {
		s += fmt.Sprintf("%s%s", v, appendStr)
	}
	return strings.TrimSuffix(s, appendStr)
}

func SubSlash(s string) []string {
	var slash string
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}

	return strings.Split(s, slash)
}
