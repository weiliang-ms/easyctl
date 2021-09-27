package exec

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	b, _ := os.ReadFile("asset/executor.yaml")
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	err := runner.Run(b, logger)
	if err != nil {
		t.Error(err)
	}
}

func TestParse(t *testing.T) {
	b, _ := os.ReadFile("asset/executor.yaml")
	executor, err := runner.ParseExecutor(b)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", executor)
}

func TestParseServerList(t *testing.T) {
	b, _ := os.ReadFile("asset/executor.yaml")
	serverList, err := runner.ParseServerList(b)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", serverList)
}
