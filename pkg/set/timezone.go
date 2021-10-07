package set

import (
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

const setTimezoneShell = "ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime"

//Timezone 设置上海时区
func Timezone(item command.OperationItem) error {
	return Config(item.B, item.Logger, setTimezoneShell)
}
