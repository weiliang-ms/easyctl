package set

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"testing"
)

func TestAddDnsScript(t *testing.T) {
	b := []byte(`
dns:
  114.114.114.114
  - 8.8.8.8
`)
	item := command.OperationItem{
		B:      b,
		Logger: logrus.New(),
	}

	assert.NotNil(t, Dns(item))

	// test invalid address
	item.B = []byte(`
dns:
  - 114.114.114.114
  - 666.555.341.11
`)
	assert.EqualError(t, Dns(item), "666.555.341.11地址非法")

	// test valid address
	assert.NotNil(t, Dns(item))
	item.B = []byte(`
dns:
  - 114.114.114.114
  - 8.8.8.8
`)
	assert.Equal(t, command.RunErr{}, Dns(item))
}
