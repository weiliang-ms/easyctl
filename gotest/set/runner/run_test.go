package runner

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	executor, err := runner.ParseExecutor(b)
	executor.Script = "mv /etc/sudoers.old /etc/sudoers"
	if err != nil {
		panic(err)
	}
	executor.ParallelRun(true)
}
