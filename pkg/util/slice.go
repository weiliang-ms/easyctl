package util

func SliceContain(s []string, element string) bool {
	for _, v := range s {
		if element == v {
			return true
		}
	}
	return false
}
