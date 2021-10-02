package util

// SliceContain 字符串切片是否含有字符串
func SliceContain(s []string, element string) bool {
	for _, v := range s {
		if element == v {
			return true
		}
	}
	return false
}
