package sys

import (
	"easyctl/util"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// https://getsubstrate.io/

const (
	LinuxGNU = "linux-gnu"
	Darwin   = "darwin"
	Freebsd  = "freebsd"
	Unknown  = "unknown system type"
)

const (
	RedhatReleaseType   = "redhat"
	RedhatOSVersionFile = "/etc/redhat-release"
)

var SystemInfoObject SystemInfo

type ServiceActionCommand struct {
	CloseFirewalld        string // 关闭防火墙
	CloseFirewalldForever string // 永久关闭防火墙
}

type OSVersion struct {
	ReleaseContent    string // release完整信息
	OSType            string
	ReleaseType       string
	MainVersionNumber string
}

type SystemInfo struct {
	Hostname      string
	Kernel        string
	RunLevel      string // 系统运行级别：0-6
	OSVersion     OSVersion
	ServiceAction ServiceActionCommand
}

func init() {
	SystemInfoObject.loadOSReleaseContent()
	SystemInfoObject.loadGnuSystemMainVersion()
	SystemInfoObject.loadRedhatManageServiceCmd()
}

// todo 待优化代码
func (system *SystemInfo) loadOSReleaseContent() {
	fmt.Println("获取操作系统版本信息...")
	// todo 优化获取os类型代码
	systemType, _ := util.ExecuteCmdAcceptResult("echo $OSTYPE")
	system.OSVersion.OSType = systemType
}

func (system *SystemInfo) loadGnuSystemMainVersion() {
	_, err := os.Stat(RedhatOSVersionFile)
	if err == nil {
		b, err := ioutil.ReadFile(RedhatOSVersionFile)
		if err == nil {
			content := string(b)
			system.OSVersion.ReleaseContent = content
			// 赋值主版本号
			system.redhatMainVersion(content)
			// 赋值release类型
			system.OSVersion.ReleaseType = RedhatReleaseType
		}
	}
}

// redhat 主版本号，如 6
func (system *SystemInfo) redhatMainVersion(releaseContent string) {
	//fmt.Println("获取redhat系版主版本号...")
	arr := strings.Split(releaseContent, " ")
	for _, v := range arr {
		//fmt.Println(v)
		matched, _ := regexp.MatchString("^[0-9].*.[0-9]$", v)
		if matched {
			number := fmt.Sprintf(strings.Split(v, ".")[0])
			system.OSVersion.MainVersionNumber = number
		}
	}
}

// 装载操作系统服务管理指令
func (system *SystemInfo) loadManageServiceCommand() {
	if system.OSVersion.OSType == LinuxGNU {
		if system.OSVersion.ReleaseType == RedhatReleaseType {
			system.loadRedhatManageServiceCmd()
		}
	} else {
		fmt.Println(errors.New("未被识别的操作系统类型"))
	}
}

// redhat系列主机管理操服务
func (system *SystemInfo) loadRedhatManageServiceCmd() {

	version := system.OSVersion.MainVersionNumber
	var closeFirewalldCmd, closeFirewalldCmdForever string

	if version == "6" {
		closeFirewalldCmd = "service" + " iptables " + "stop"
		closeFirewalldCmdForever = ""
	} else if version == "7" {
		closeFirewalldCmd = "systemctl" + " stop " + "firewalld"
		closeFirewalldCmdForever = "systemctl" + " disable " + "firewalld"
	}

	// todo 组装
	system.ServiceAction.CloseFirewalld = closeFirewalldCmd
	system.ServiceAction.CloseFirewalldForever = closeFirewalldCmdForever
}

func (system *SystemInfo) loadRunLevel() {
	level, err := util.ExecuteCmdAcceptResult("runlevel")
	if err == nil {
		system.RunLevel = level
	}
}

func (system *SystemInfo) kernelVersion() {
	cmd := "uname -r"
	kernel, err := util.ExecuteCmdAcceptResult(cmd)
	if err == nil {
		system.Kernel = kernel
	}
}

func PrintSystemInfo() {
	fmt.Printf("###   当前系统信息   ####\n\n"+
		"操作系统版本：%s",
		SystemInfoObject.OSVersion.ReleaseContent)
}
