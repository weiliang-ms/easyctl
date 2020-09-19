package sys

import (
	"easyctl/util"
	"fmt"
)

const closeSelinuxFailure = "关闭selinux失败..."

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
		closeSelinuxSuccessful += "\n执行 init 6 重启主机命令保证selinux永久关闭生效"
		util.PrintSuccessfulMsg(fmt.Sprintf("%s", closeSelinuxSuccessful))
	} else if err == nil && !forever {
		util.PrintSuccessfulMsg(fmt.Sprintf("%s", closeSelinuxSuccessful))
	}
}
