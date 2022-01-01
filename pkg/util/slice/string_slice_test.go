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

func TestStringSliceFilter(t *testing.T) {
	source1 := []string{"1", "", "2", "", "", "3"}
	expect1 := []string{"1", "2", "3"}
	assert.Equal(t, expect1, StringSliceFilter(source1, ""))

	source2 := []string{"1", "c", "2", "b", "c", "3"}
	expect2 := []string{"1", "2", "b", "3"}
	assert.Equal(t, expect2, StringSliceFilter(source2, "c"))
}
