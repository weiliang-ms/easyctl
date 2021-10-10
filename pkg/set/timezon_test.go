package set

import (
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"testing"
)

func TestTimezone(t *testing.T) {
	assert.Equal(t, command.RunErr{}, Timezone(command.OperationItem{}))
}
