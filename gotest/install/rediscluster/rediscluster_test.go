package ulimit

import (
	_ "embed"
	"github.com/weiliang-ms/easyctl/pkg/install"
	"os"
	"testing"
)

func TestUlimit(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	err := install.RedisCluster(b, false)
	if err != nil {
		panic(err)
	}

}
