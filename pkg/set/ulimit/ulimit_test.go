package ulimit

import (
	_ "embed"
	"testing"
)

func TestUlimit(t *testing.T) {
	err := Config("config.yaml", 1)
	if err != nil {
		panic(err)
	}
}
