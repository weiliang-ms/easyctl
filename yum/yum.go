package yum

import (
	"easyctl/constant"
	"easyctl/util"
	"fmt"
)

func Detection(packageName string, origin bool, instance util.SSHInstance) (re bool) {
	if !origin {
		fmt.Printf("%s 检测%s是否已安装\n", util.PrintOrange(constant.Check), packageName)
		return util.ExecuteIgnoreStd(fmt.Sprintf("rpm -qa|grep %s", packageName))
	} else {
		fmt.Printf("%s 远程%s安装%s\n", util.PrintOrange(constant.Install), instance.Host, packageName)
		return util.ExecuteIgnoreStd(fmt.Sprintf("yum install -y %s", packageName))
	}
}

func Install(packageName string) bool {
	fmt.Printf("%s 安装%s\n", util.PrintOrange(constant.Install), packageName)
	return util.ExecuteIgnoreStd(fmt.Sprintf("yum install -y %s", packageName))
}
