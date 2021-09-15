package util

import (
	"fmt"
	"strconv"
	"strings"
)

func KConvert(value string) string {
	b, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(err)
	}
	if b > 1024 {
		v := b / 1024
		if v > 1024 {
			return fmt.Sprintf("%sG", ChangeNumber(v/1024, 1))
		}
		return fmt.Sprintf("%sM", ChangeNumber(v, 1))
	}
	return fmt.Sprintf("%sK", ChangeNumber(b, 1))
}

func ChangeNumber(f float64, m int) string {
	n := strconv.FormatFloat(f, 'f', -1, 64)
	if n == "" {
		return ""
	}
	if m >= len(n) {
		return n
	}
	newn := strings.Split(n, ".")
	if len(newn) < 2 || m >= len(newn[1]) {
		return n
	}
	return newn[0] + "." + newn[1][:m]
}
