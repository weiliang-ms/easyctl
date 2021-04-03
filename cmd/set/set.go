package set

import (
	"errors"
	"fmt"
)

// todo
var (
	CheckIDRsaFileCmd    = fmt.Sprintf("[ -e $HOME/.ssh/%s ]", idRsa)
	CheckIDRsaPubFileCmd = fmt.Sprintf("[ -e $HOME/.ssh/%s ]", idRsaPub)

	HostsFilePath string // 主机互信hosts文件路径

	missingParameterErr = errors.New("缺少参数...")
)

//// set 命令合法参数
var setValidArgs = []string{"dns", "yum", "hostname", "timezone", "password-less"}
