package errors

import "fmt"

// NumNotEqualErr 数量不匹配error
func NumNotEqualErr(msg string, expect, acture int) error {
	return fmt.Errorf("%s数量非法 expect num: %d but get: %d", msg, expect, acture)
}

// FileNotFoundErr 数量不匹配error
func FileNotFoundErr(filepath string) error {
	return fmt.Errorf("%s 非法路径", filepath)
}
