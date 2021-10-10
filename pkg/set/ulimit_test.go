package set

import (
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"testing"
)

func TestUlimit(t *testing.T) {
	err := Ulimit(command.OperationItem{
		B:          nil,
		Logger:     nil,
		OptionFunc: nil,
	})

	assert.Equal(t, command.RunErr{}, err)
}
