package set

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

const setTimezoneShell = "ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime"

//Timezone 设置上海时区
func Timezone(item command.OperationItem) error {
	return runner.RemoteRun(item.B, item.Logger, setTimezoneShell)
}
