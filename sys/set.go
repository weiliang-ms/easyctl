package sys

import (
	"easycfg/resources"
	"easycfg/util"
	"fmt"
	"log"
	"os"
)

var aliBaseEL7WriteError = "阿里云base镜像源配置失败..."
var aliEpelEL7WriteError = "阿里云epel镜像源配置失败..."
var setAliMirrorSuccessful = "阿里云镜像源配置成功..."

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

	_, baseWriteErr := baseRepoFile.Write([]byte(resources.CentOS7AliBaseYUMContent))
	if baseWriteErr != nil {
		fmt.Println(baseWriteErr.Error())
		fmt.Println("[failed] " + aliBaseEL7WriteError)
	}

	fmt.Printf("[create] 创建epel-ali.repo文件...\n")
	epelRepoFile, err := os.OpenFile("/etc/yum.repos.d/epel-ali.repo", os.O_WRONLY|os.O_CREATE, 0666)
	defer epelRepoFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	_, epelWriteErr := epelRepoFile.Write([]byte(resources.CentOS7AliEpelYUMContent))
	if epelWriteErr != nil {
		fmt.Println(epelWriteErr.Error())
		fmt.Println("[failed] " + aliEpelEL7WriteError)
	}

	cleanYUMCacheCmd := "yum clean all"
	fmt.Printf("[clean] 清除yum缓存...\n")
	util.ExecuteCmd(cleanYUMCacheCmd)

	fmt.Println("[successful] " + setAliMirrorSuccessful)
}

func SetHostname(name string) {
	// todo
	fmt.Println("校验hostname格式逻辑...")
	fmt.Println("[hostname]配置hostname...")
	cmd := fmt.Sprintf("hostnamectl --static set-hostname %s", name)
	err, _ := util.ExecuteCmd(cmd)
	if err != nil {
		log.Println(err.Error())
		util.PrintFailureMsg("[failed] 配置hostname失败...")
	} else {
		util.PrintSuccessfulMsg("[success] 配置hostname成功...")
	}

	fmt.Println("[host]配置host解析...")
	util.ExecuteCmd(fmt.Sprintf("echo \"127.0.0.1 %s\"", name))
}
