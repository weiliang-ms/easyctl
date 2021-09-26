package ulimit

import (
	_ "embed"
	"github.com/weiliang-ms/easyctl/pkg/set"
	"os"
	"testing"
)

func TestUlimit(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	err := set.Ulimit(b, true)
	if err != nil {
		panic(err)
	}

}