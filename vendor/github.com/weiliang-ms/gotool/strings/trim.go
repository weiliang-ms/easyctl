package strings

import (
	"strings"
	"unicode"
)

func TrimNumSuffix(str string) string {

	end := str[len(str)-1]
	r := rune(end)
	endChar := str[len(str)-1:]

	if unicode.IsNumber(r) {
		return strings.TrimSuffix(str, endChar)
	}

	return str
}
