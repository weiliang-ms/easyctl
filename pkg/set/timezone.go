package set

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

const setTimezoneShell = "ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime"

//Timezone 设置上海时区
func Timezone(item command.OperationItem) command.RunErr {
	return runner.RemoteRun(runner.RemoteRunItem{
		ManifestContent:     item.B,
		Logger:              item.Logger,
		Cmd:                 setTimezoneShell,
		RecordErrServerList: false,
	})
}
