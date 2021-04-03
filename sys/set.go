package sys

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/constant"
	"github.com/weiliang-ms/easyctl/util"
	"log"
	"os"
)

const (
	aliBaseEL7WriteErrMsg = "阿里云base镜像源配置失败..."
	localWriteErrMsg      = "local.repo文件写失败..."

	nginxRepoFileWriteErrMsg = "nginx.repo文件写失败..."

	aliEpelEL7WriteErrMsg = "阿里云epel镜像源配置失败..."

	setAliMirrorSuccessful   = "阿里云镜像源配置成功..."
	setLocalMirrorSuccessful = "本地镜像源配置成功..."
	setNginxMirrorSuccessful = "nginx镜像源配置成功..."
)

func SetDNS(dnsAddress string) (err error, result string) {

	cmd := "sed -i \"/nameserver " + dnsAddress + "/d\" /etc/resolv.conf;" +
		"echo \"nameserver " + dnsAddress + "\" >> /etc/resolv.conf\n"

	fmt.Printf("[check] 检测dns地址：%s合法性...\n", dnsAddress)
	err, result = util.CheckIP(dnsAddress)
	if err != nil {
		return err, result
	}
	shellErr, shellResult := util.ExecuteCmd(cmd)
	return shellErr, shellResult
}

func SetNginxMirror() {

	// 写nginx.repo文件
	fmt.Printf("[yum] 创建nginx.repo文件...\n")
	nginxRepoFile, err := os.OpenFile("/etc/yum.repos.d/nginx.repo", os.O_WRONLY|os.O_CREATE, 0666)
	defer nginxRepoFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	_, nginxRepoFileWriteErr := nginxRepoFile.Write([]byte(constant.CentOSNginxMirrorContent))
	if nginxRepoFileWriteErr != nil {
		fmt.Println(nginxRepoFileWriteErr.Error())
		fmt.Println("[failed] " + nginxRepoFileWriteErrMsg)
	}

	cleanYUMCacheCmd := "yum clean all"
	fmt.Printf("[clean] 清除yum缓存...\n")
	util.ExecuteCmd(cleanYUMCacheCmd)

	fmt.Println("[successful] " + setNginxMirrorSuccessful)
}

func SetHostname(name string) {
	// todo 校验hostname格式逻辑
	fmt.Println("[hostname]配置hostname...")

	cmd := fmt.Sprintf("sed -i '/HOSTNAME/d' /etc/sysconfig/network;"+
		"echo \"HOSTNAME=%s\" >> /etc/sysconfig/network;"+
		"sysctl kernel.hostname=%s", name, name)

	err, _ := util.ExecuteCmd(cmd)

	if err != nil {
		log.Println(err.Error())
		util.PrintFailureMsg("[failed] 配置hostname失败...")
	} else {
		util.PrintSuccessfulMsg("[success] 配置hostname成功...")
	}

	fmt.Println("[host]配置host解析...")
	util.ExecuteCmd(fmt.Sprintf("echo \"127.0.0.1 %s\" >> /etc/hosts", name))
}
