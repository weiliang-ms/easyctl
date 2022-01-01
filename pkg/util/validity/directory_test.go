package validity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataPath(t *testing.T) {
	var err error

	err = DataPath("./ddd")
	assert.Equal(t, relativePathErr{Msg: "./ddd"}, err)

	err = DataPath("ddd/aaa/fff")
	assert.Equal(t, relativePathErr{Msg: "ddd/aaa/fff"}, err)

	err = DataPath("/proc/ddd")
	assert.Equal(t, sensitivePathErr{Msg: "/proc"}, err)

	err = DataPath("/boot/ddd")
	assert.Equal(t, sensitivePathErr{Msg: "/boot"}, err)

	err = DataPath("/bin/ddd")
	assert.Equal(t, sensitivePathErr{Msg: "/bin"}, err)

	err = DataPath("/sbin/ddd")
	assert.Equal(t, sensitivePathErr{Msg: "/sbin"}, err)

	err = DataPath("/ddd")
	assert.Equal(t, nil, err)
}

func TestErrorContent(t *testing.T) {
	relativePathErr := relativePathErr{Msg: "bin"}
	assert.Equal(t, "bin 路径非法，不是绝对路径", relativePathErr.Error())

	sensitivePathErr := sensitivePathErr{Msg: "/bin"}
	assert.Equal(t, "路径非法，不允许在/bin目录下", sensitivePathErr.Error())
}