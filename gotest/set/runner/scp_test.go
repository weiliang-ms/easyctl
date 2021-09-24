package runner

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"os"
	"testing"
)

const (
	scpSrc = "../../../asset/iperf3-3.1.7-2.el7.x86_64.rpm"
	scpDst = "/tmp/iperf3-3.1.7-2.el7.x86_64.rpm"
)

func TestScp(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	servers, err := runner.ParseServerList(b)
	if err != nil {
		panic(err)
	}

	if len(servers) > 0 {
		scpErr := servers[0].Scp(scpSrc, scpDst, 0666, true)
		if scpErr != nil {
			panic(scpErr)
		}
	}

}

func TestParallelScp(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	servers, err := runner.ParseServerList(b)
	if err != nil {
		panic(err)
	}

	_ = runner.ParallelScp(servers, scpSrc, scpDst, 0755)

}
