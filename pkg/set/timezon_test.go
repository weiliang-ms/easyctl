package set

import (
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"testing"
)

func TestTimezone(t *testing.T) {
	assert.Nil(t, Timezone(command.OperationItem{}))
}
