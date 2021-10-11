package deny

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

// Firewall 关闭防火墙
func Firewall(item command.OperationItem) command.RunErr {
	return runner.RemoteRun(item.B, item.Logger, disableFirewallShell)
}
