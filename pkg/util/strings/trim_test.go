package strings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrimNumSuffix(t *testing.T) {

	assert.Equal(t, "/dev/vda", TrimNumSuffix("/dev/vda1"))
	assert.Equal(t, "/dev/vdad", TrimNumSuffix("/dev/vdad"))
}
