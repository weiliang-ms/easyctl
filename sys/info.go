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
	start               = "start"
	stop                = "stop"
	disable             = "disable"
	on                  = "on"
	off                 = "off"
	service             = "service"
	enable              = "enable"
	chkconfig           = "chkconfig"
	systemctl           = "systemctl"
)

var SystemInfoObject SystemInfo

type ServiceActionCommand struct {
	CloseFirewalld        string // 关闭防火墙
	CloseFirewalldForever string // 永久关闭防火墙
	StartDocker           string // 开启docker服务
	StartDockerForever    string // 永久开启docker服务
	StartNginx            string // 开启nginx服务
	StartNginxForever     string // 永久开启nginx服务
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
	version := SystemInfoObject.OSVersion.MainVersionNumber
	// todo 组装
	system.ServiceAction.CloseFirewalld = packageRedhatService(version, "firewalld", stop, false)
	system.ServiceAction.CloseFirewalldForever = packageRedhatService(version, "firewalld", stop, true)
	system.ServiceAction.StartDocker = packageRedhatService(version, "docker", start, false)
	system.ServiceAction.StartDockerForever = packageRedhatService(version, "docker", start, true)
	system.ServiceAction.StartNginx = packageRedhatService(version, "nginx", start, false)
	system.ServiceAction.StartNginxForever = packageRedhatService(version, "nginx", start, true)
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
