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
}
