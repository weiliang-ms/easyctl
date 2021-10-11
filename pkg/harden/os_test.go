package harden

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestOS(t *testing.T) {
	err := OS(command.OperationItem{Logger: logrus.New()})
	assert.Equal(t, command.RunErr{}, err)

	// 模式异常
	err = OS(command.OperationItem{Logger: logrus.New(), B: []byte(`
server:
   host: 1.1.1.1
`)})

	_, ok := err.Err.(*yaml.TypeError)
	assert.Equal(t, true, ok)
}
