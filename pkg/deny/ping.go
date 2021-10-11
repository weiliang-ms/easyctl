package deny

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

// Ping 禁ping
func Ping(item command.OperationItem) command.RunErr {
	return command.RunErr{Err: runner.RemoteRun(item.B, item.Logger, DenyPingShell)}
}
