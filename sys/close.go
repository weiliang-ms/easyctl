package sys

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/util"
)

const closeSelinuxFailure = "关闭selinux失败..."

// todo 待优化变量
var closeSelinuxSuccessful = "关闭selinux成功..."

func CloseSeLinux(forever bool) {
	cmd := "setenforce 0;"
	if forever {
		cmd += "sed -i \"s#SELINUX=enforcing#SELINUX=disabled#g\" /etc/selinux/config"
	}

	err, _ := util.ExecuteCmd(cmd)

	if err != nil {
		util.PrintFailureMsg(fmt.Sprintf("%s\n%s", closeSelinuxFailure, err.Error()))
	} else if err == nil && forever {
		// todo 待优化
		closeSelinuxSuccessful += "\n执行 init 6 重启主机命令保证selinux永久关闭生效"
		util.PrintSuccessfulMsg(fmt.Sprintf("%s", closeSelinuxSuccessful))
	} else if err == nil && !forever {
		util.PrintSuccessfulMsg(fmt.Sprintf("%s", closeSelinuxSuccessful))
	}
}

// 关闭防火墙
func CloseFirewalld(forever bool) {

	var cmd string
	switch forever {
	case true:
		cmd = SystemInfoObject.ServiceAction.CloseFirewalldForever
	case false:
		cmd = SystemInfoObject.ServiceAction.CloseFirewalld
	}

	util.ExecuteCmdAcceptResult(cmd)
}
