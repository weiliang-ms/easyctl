package errors

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"runtime"
)

// NumNotEqualErr 数量不匹配error
func NumNotEqualErr(msg string, expect, acture int) error {
	return fmt.Errorf("%s数量非法 expect num: %d but get: %d", msg, expect, acture)
}

// FileNotFoundErr 数量不匹配error
func FileNotFoundErr(filepath string) error {
	return fmt.Errorf("%s 非法路径", filepath)
}

// IgnoreErrorFromCaller 忽略来自指定调用者的异常（测试用例）
func IgnoreErrorFromCaller(skip int, callerName string, err *error) {
	pc, _, _, _ := runtime.Caller(skip)
	name := runtime.FuncForPC(pc).Name()
	if name == callerName && err != nil {
		*err = nil
	}
}

/*
	condition == true  -> return err
	condition == false -> return nil
*/
func FalseConditionErr(condition bool, msg string) error {
	if !condition {
		return nil
	}
	return fmt.Errorf(fmt.Sprintf("%s.", msg))
}

// IsTestCaller 判断调用者是否为test函数
func IsTestCaller(skip int) bool {
	pc, _, _, _ := runtime.Caller(skip)
	name := runtime.FuncForPC(pc).Name()
	if name == constant.TRunnerCaller {
		return true
	}
	return false
}

func IsCaller(skip int, callerName string) bool {
	pc, _, _, _ := runtime.Caller(skip)
	name := runtime.FuncForPC(pc).Name()
	if name == callerName {
		return true
	}
	return false
}
