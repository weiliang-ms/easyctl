package deny

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"testing"
)

func TestFirewall(t *testing.T) {
	err := Firewall(command.OperationItem{
		B:          nil,
		Logger:     logrus.New(),
		OptionFunc: nil,
	})

	assert.Equal(t, nil, err)
}
