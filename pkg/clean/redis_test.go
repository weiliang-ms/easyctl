package clean

import (
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"testing"
)

func TestRedis(t *testing.T) {
	runErr := Redis(command.OperationItem{})
	assert.Nil(t, runErr.Err)
}
