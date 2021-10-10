package set

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"testing"
)

func TestMakeKeyPairScript(t *testing.T) {
	content, err := MakeKeyPairScript(PasswordLessTmpl)
	assert.Nil(t, err)
	assert.NotEqual(t, "", content)
}

func TestPasswordLess(t *testing.T) {
	err := PasswordLess(command.OperationItem{
		B:          nil,
		Logger:     logrus.New(),
		OptionFunc: nil,
	})

	assert.Equal(t, command.RunErr{}, err)
}
