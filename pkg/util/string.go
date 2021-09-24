package util

import (
	"fmt"
	"os"
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

func SubFileName(s string) string {

	f, err := os.Stat(s)
	if err != nil {
		panic(err)
	}
	return f.Name()
}
