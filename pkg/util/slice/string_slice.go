package slice

func StringSliceRemove(sup []string, cut []string) []string {
	var re []string
	for _, v := range sup {
		if !StringSliceContain(cut, v) {
			re = append(re, v)
		}
	}
	return re
}

// StringSliceContain todo 优化判断方式
func StringSliceContain(sup []string, s string) bool {
	for _, v := range sup {
		if v == s {
			return true
		}
	}
	return false
}