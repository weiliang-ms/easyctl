package runner

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"os"
	"testing"
)

const (
	scpSrc = "C:\\Users\\weiliang\\Desktop\\sentinel-1.8.1.zip"
	scpDst = "/opt/sentinel-1.8.1.zip"
)

//func TestScp(t *testing.T) {
//	b, readErr := os.ReadFile("../../../asset/config.yaml")
//	if readErr != nil {
//		panic(readErr)
//	}
//
//	servers, err := runner.ParseServerList(b)
//	if err != nil {
//		panic(err)
//	}
//
//	if len(servers) > 0 {
//		scpErr := servers[0].Scp(scpSrc, scpDst, 0666, true)
//		if scpErr != nil {
//			panic(scpErr)
//		}
//	}
//
//}

func TestParallelScp(t *testing.T) {
	b, readErr := os.ReadFile("../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	servers, err := runner.ParseServerList(b)
	if err != nil {
		panic(err)
	}

	_ = runner.ParallelScp(servers, scpSrc, scpDst, 0755)

}
