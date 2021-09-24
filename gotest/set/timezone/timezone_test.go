package ulimit

import (
	_ "embed"
	"github.com/weiliang-ms/easyctl/pkg/set"
	"os"
	"testing"
)

func TestTimezone(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	err := set.Timezone(b, true)
	if err != nil {
		panic(err)
	}
}
