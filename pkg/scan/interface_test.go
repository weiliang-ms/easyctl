package scan

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"testing"
)

func Test(t *testing.T) {

	s := runner.ServerInternal{}
	r := OsExecutor{}

	_, _ = r.GetHostName(s, mockLogger)
	_, _ = r.GetCPUInfo(s, mockLogger)
	_, _ = r.GetCPULoadAverage(s, mockLogger)
	_, _ = r.GetKernelVersion(s, mockLogger)
	_, _ = r.GetSystemVersion(s, mockLogger)
	_, _ = r.GetMountPointInfo(s, mockLogger)
	_, _ = r.GetMemoryInfo(s, mockLogger)

	r.DoRequest(runner.DoRequestItem{
		S: s,
		R: runner.RunItem{
			Logger: mockLogger,
			Cmd:    "ddd",
		},
		Mock: true,
	})

}
