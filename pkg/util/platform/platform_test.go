package platform

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSlash(t *testing.T) {
	require.Equal(t, "/", Slash("linux"))
	require.Equal(t, "\\", Slash("windows"))
	require.Equal(t, "/", Slash("other"))
	require.Equal(t, "/", Slash(""))
}
