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

func SetAliYUM() {

	// 备份repo文件
	cmd := "mkdir -p /etc/yum.repos.d/`date +%Y%m%d`" + ";" +
		"mv /etc/yum.repos.d/*.repo /etc/yum.repos.d/`date +%Y%m%d` -f"

	fmt.Printf("[bakup] 备份历史repo文件...\n")
	util.ExecuteCmd(cmd)

	// 写入base文件
	fmt.Printf("[create] 创建base-ali.repo文件...\n")
	baseRepoFile, err := os.OpenFile("/etc/yum.repos.d/base-ali.repo", os.O_WRONLY|os.O_CREATE, 0666)
	defer baseRepoFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	_, baseWriteErr := baseRepoFile.Write([]byte(constant.CentOSAliBaseYUMContent))
	if baseWriteErr != nil {
		fmt.Println(baseWriteErr.Error())
		fmt.Println("[failed] " + aliBaseEL7WriteErrMsg)
	}

	fmt.Printf("[create] 创建epel-ali.repo文件...\n")
	epelRepoFile, err := os.OpenFile("/etc/yum.repos.d/epel-ali.repo", os.O_WRONLY|os.O_CREATE, 0666)
	defer epelRepoFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	_, epelWriteErr := epelRepoFile.Write([]byte(constant.CentOSAliEpelYUMContent))
	if epelWriteErr != nil {
		fmt.Println(epelWriteErr.Error())
		fmt.Println("[failed] " + aliEpelEL7WriteErrMsg)
	}

	cleanYUMCacheCmd := "yum clean all"
	fmt.Printf("[clean] 清除yum缓存...\n")
	util.ExecuteCmd(cleanYUMCacheCmd)

	fmt.Println("[successful] " + setAliMirrorSuccessful)
}

func SetLocalYUM() {

	// 备份repo文件
	cmd := "mkdir -p /etc/yum.repos.d/`date +%Y%m%d`" + ";" +
		"mv /etc/yum.repos.d/*.repo /etc/yum.repos.d/`date +%Y%m%d` -f"

	fmt.Printf("[bakup] 备份历史repo文件...\n")
	util.ExecuteCmd(cmd)

	// 写local.repo文件
	fmt.Printf("[create] 创建local.repo文件...\n")
	localRepoFile, err := os.OpenFile("/etc/yum.repos.d/local.repo", os.O_WRONLY|os.O_CREATE, 0666)
	defer localRepoFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	_, localWriteErr := localRepoFile.Write([]byte(constant.CentOSLocalYUMContent))
	if localWriteErr != nil {
		fmt.Println(localWriteErr.Error())
		fmt.Println("[failed] " + localWriteErrMsg)
	}

	cleanYUMCacheCmd := "yum clean all"
	fmt.Printf("[clean] 清除yum缓存...\n")
	util.ExecuteCmd(cleanYUMCacheCmd)

	fmt.Println("[successful] " + setLocalMirrorSuccessful)
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

func SetTimeZone() {
	fmt.Println("[timezone]配置时区为上海...")
	util.ExecuteCmd("\\cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime -R")
}
