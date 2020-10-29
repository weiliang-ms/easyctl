package util

import (
	"strings"
)

func CutCharacter(str string, cutCharacters []string) (s string) {
	for _, c := range cutCharacters {
		s = strings.Replace(str, c, "", -1)
	}
	return strings.TrimSpace(s)
}
