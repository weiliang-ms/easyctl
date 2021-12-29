package clean

import (
	"github.com/weiliang-ms/easyctl/pkg/install"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

// Redis 清理redis服务及文件
func Redis(item command.OperationItem) command.RunErr {
	return runner.RemoteRun(runner.RemoteRunItem{
		B:                   item.B,
		Logger:              item.Logger,
		Cmd:                 install.PruneRedisShell,
		RecordErrServerList: false,
	})
}
