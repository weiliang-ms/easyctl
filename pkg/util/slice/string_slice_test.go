package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringSliceContain(t *testing.T) {
	slice := []string{"aaa", "bbb", "ccc"}
	s := "bbb"
	assert.Equal(t, true, StringSliceContain(slice, s))
}

func TestStringSliceRemove(t *testing.T) {
	slice := []string{"aaa", "bbb", "ccc", "ddd"}
	subSlice := []string{"bbb", "ccc"}
	assert.Equal(t, []string{"aaa", "ddd"}, StringSliceRemove(slice, subSlice))
}
