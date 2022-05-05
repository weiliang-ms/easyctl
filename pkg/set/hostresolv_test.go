package set

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"runtime"
	"testing"
)

func mockHostResolveFuncWithNilData(b []byte, logger *logrus.Logger, cmd string) ([]runner.ShellResult, error) {
	return []runner.ShellResult{}, fmt.Errorf("ddd")
}

func mockHostResolveFuncWithData(b []byte, logger *logrus.Logger, cmd string) ([]runner.ShellResult, error) {
	return []runner.ShellResult{
		{
			Host:   "1.1.1.1",
			StdOut: "server-A",
		},
		{
			Host:   "1.1.1.2",
			StdOut: "server-B",
		},
		{
			Host:   "1.1.1.3",
			StdOut: "localhost",
		},
	}, nil
}

func TestHostResolve(t *testing.T) {
	item := command.OperationItem{}
	item.B = []byte(`
server:
  - host: 114.114.114.114
    username: "root"
    password: 123
    port: 22
`)

	item.Logger = logrus.New()
	options := make(map[string]interface{})
	options[GetHostResolveFunc] = GetHostResolve
	item.OptionFunc = options
	err := HostResolve(item)
	_, ok := err.Err.(runtime.Error)
	assert.Equal(t, false, ok)

	//// test mock with nil data
	//options[GetHostResolveFunc] = mockHostResolveFuncWithNilData
	//item.OptionFunc = options
	//err = HostResolve(item)
	//assert.EqualError(t, err, "ddd")
	//
	//// test mock with data
	//options[GetHostResolveFunc] = mockHostResolveFuncWithData
	//item.OptionFunc = options
	//item.B = []byte{}
	//err = HostResolve(item)
	//assert.Equal(t, command.RunErr{}, err)
	//
	//// test bad function
	//options[GetHostResolveFunc] = HostResolve
	//item.OptionFunc = options
	//item.B = []byte{}
	//err = HostResolve(item)
	//assert.EqualError(t, err, "入参：getHostResolveFunc 非法")
}
