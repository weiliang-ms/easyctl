package errors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"testing"
)

func TestIsTestCaller(t *testing.T) {
	assert.Equal(t, false, IsTestCaller(1))
	assert.Equal(t, true, IsTestCaller(2))
	assert.Equal(t, false, IsTestCaller(3))
}

func TestIgnoreErrorFromCaller(t *testing.T) {
	err := errors.New("test error")
	IgnoreErrorFromCaller(2, constant.TRunnerCaller, &err)
	assert.Nil(t, err)
}

func TestFileNotFoundErr(t *testing.T) {
	assert.EqualError(t, FileNotFoundErr("1.txt"), "1.txt 非法路径")
}

func TestNumNotEqualErr(t *testing.T) {
	assert.EqualError(t, NumNotEqualErr("ddd", 1, 3), "ddd数量非法 expect num: 1 but get: 3")
}

func TestIsCaller(t *testing.T) {
	assert.Equal(t, true, IsCaller(1, "github.com/weiliang-ms/easyctl/pkg/util/errors.TestIsCaller"))
}
