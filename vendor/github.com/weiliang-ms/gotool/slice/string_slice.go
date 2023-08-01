package slice

// StringSliceRemove 字符串切片过滤字符串数组中的元素
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
func StringSliceContain(sup []string, element string) bool {
	for _, v := range sup {
		if v == element {
			return true
		}
	}
	return false
}

// StringSliceFilter 过滤字符串数组元素
func StringSliceFilter(s []string, filterChar string) []string {
	var r []string
	for _, v := range s {
		if v != filterChar {
			r = append(r, v)
		}
	}
	return r
}

// StringSliceAppend 字符串数组拼接
func StringSliceAppend(s []string, sub []string) []string {
	for _, v := range sub {
		s = append(s, v)
	}
	return s
}
