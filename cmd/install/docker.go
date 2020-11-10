package install

import (
	"easyctl/constant"
	"easyctl/sys"
	"easyctl/util"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type docker struct {
	serverList      []util.Server
	serviceFilePath string
	serviceContent  string
}

const (
	parseDockerServerList     = "解析安装docker宿主机列表信息"
	decompressDockerPackage   = "解压docker安装包文件"
	reloadUnitService         = "重载系统service unit配置文件"
	enableDcokerService       = "配置docker为系统服务"
	startDockerService        = "启动docker服务"
	statusDockerService       = "docker服务状态"
	disableSelinux            = "关闭selinux"
	checkSelinuxStatus        = "检测selinux状态"
	checkKernelVersion        = "检测内核版本"
	afterModifySelinuxContent = "执行root重启主机，以保证selinux关闭，随后重新执行安装命令进行安装"
	decompressDockerCmd       = "tar zxvf docker-*.tgz;mv docker/* /usr/bin/"
	scpDockerPackage          = "拷贝docker安装包"
)

func (docker *docker) Install() {
	if dockerOffline {
		docker.offline()
	} else {
		fmt.Println("在线安装...")
	}
}

func (docker *docker) offline() {
	// 解析参数
	docker.parseFlags()

	// 赋值
	docker.setDefaultValue()

	// 解析安装docker宿主机列表
	docker.getServerList()

	// 拷贝安装文件
	docker.scp()

	// 判断内核版本
	docker.checkKernel()

	// 解压
	docker.decompress()

	// 检测selinux
	docker.checkSelinux()

	// 写unit服务配置
	docker.writeService()

	docker.reload()

	docker.enable()

	docker.start()

	docker.status()

}

func (docker *docker) parseFlags() {
	if dockerOffline && dockerPackageFile == "" {
		log.Fatal("离线模式下，需通过--file= 参数指定离线安装包路径...")
	}
}

func (docker *docker) setDefaultValue() {
	docker.serviceFilePath = constant.Redhat7DockerServiceFilePath
	docker.serviceContent = constant.Redhat7DockerServiceContent
}

// docker server list赋值
func (docker *docker) getServerList() {
	if dockerServerListFile != "" {
		if _, err := os.Stat(dockerServerListFile); err == nil {
			docker.banner(constant.Parse, parseDockerServerList)
			docker.serverList = util.ParseServerList(dockerServerListFile).DockerServerList
		} else {
			log.Fatal(err.Error())
		}
	}
}

// 解压赋值
func (docker *docker) decompress() {
	cmd := fmt.Sprintf("%s && %s", constant.RootDetectionCmd, decompressDockerCmd)
	docker.shell(cmd, decompressDockerPackage).Shell()
}

// 拷贝
func (docker *docker) scp() {
	// 拷贝docker安装包
	if len(docker.serverList) != 0 {
		util.ScpHome(docker.returnBanner(scpDockerPackage), dockerPackageFile, docker.serverList)
	}
	time.Sleep(3 * time.Second)
}

// 检测内核版本
func (docker *docker) checkKernel() {
	if len(docker.serverList) == 0 {
		docker.checkLocalKernelVersion()
	} else {
		docker.checkRemoteKernelVersion()
	}
}

func (docker *docker) checkLocalKernelVersion() {
	util.PrintDirectBanner([]string{constant.LoopbackAddress}, checkKernelVersion)
	if sys.KernelMainVersion() < 3 {
		log.Fatal("内核版本小于3,不满足安装条件...")
	}
}

func (docker *docker) checkRemoteKernelVersion() {
	for _, v := range docker.serverList {
		util.PrintDirectBanner([]string{v.Host}, checkKernelVersion)
		version := v.RemoteShellReturnStd(constant.KernelVersionCmd)
		mVersion, err := strconv.Atoi(strings.TrimSuffix(version, "\n"))
		if err != nil {
			log.Fatal(err.Error())
		}
		if mVersion < 3 {
			log.Fatal("内核版本小于3,不满足安装条件...")
		}
	}
}

// 配置系统服务
func (docker *docker) writeService() {
	util.WriteFile(docker.serviceFilePath, []byte(docker.serviceContent), docker.serverList)
}

// 载入service文件
func (docker *docker) reload() {
	docker.shell(constant.DeamonReloadCmd, reloadUnitService).Shell()
}

// 开机自启动
func (docker *docker) enable() {
	docker.shell(constant.DockerEnableCmd, enableDcokerService).Shell()
}

// 检测selinux状态
func (docker *docker) checkSelinux() {
	if len(docker.serverList) == 0 {
		docker.closeLocalSelinux()
	} else {
		for _, v := range docker.serverList {
			docker.closeRemoteSelinux(v)
		}
	}
}

// 关闭本地selinux
func (docker *docker) closeLocalSelinux() {
	status := docker.localSelinuxStatus()
	if strings.TrimSuffix(status, "\n") == constant.Enabled {
		util.LocalShell(disableSelinux, constant.DisableSelinuxCmd)
		log.Fatal(afterModifySelinuxContent)
	}
}

func (docker *docker) closeRemoteSelinux(server util.Server) {
	status := docker.remoteSelinuxStatus(server)
	if strings.TrimSuffix(status, "\n") == constant.Enabled {
		server.RemoteShellReturnStd(constant.DisableSelinuxCmd)
		log.Fatal(fmt.Sprintf("[%s]%s", server.Host, afterModifySelinuxContent))
	}
}

func (docker *docker) localSelinuxStatus() string {
	return util.LocalShell(checkSelinuxStatus, constant.SelinuxStatusCmd)
}

func (docker *docker) remoteSelinuxStatus(server util.Server) string {
	return util.RemoteShellAcceptResult(checkSelinuxStatus, constant.SelinuxStatusCmd, server)
}

func (docker *docker) start() {
	docker.shell(constant.DockerStartCmd, startDockerService).Shell()
}

func (docker *docker) status() {
	docker.shell(constant.DockerStatusCmd, statusDockerService).ShellPrintStdout()
}

// 在线安装docker
func installDockerOnline() {
	fmt.Println("检测内核...")
	if !sys.AccessAliMirrors() {
		panic(netConnectErr)
	}

	sys.SetAliYUM()
	install := "yum -y redis yum-utils device-mapper-persistent-data lvm2;" +
		"yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo;" +
		"yum makecache fast;" +
		"yum -y redis docker-ce"

	//
	util.ExecuteCmdAcceptResult(install)

	sys.CloseSeLinux(true)
	fmt.Println("[docker]启动docker...")
	startRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartDocker)
	fmt.Println("[docker]设置docker开机自启动...")
	enableRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartDockerForever)
	if startRe == nil && enableRe == nil {
		util.PrintSuccessfulMsg("docker安装成功...")
	}
}

func (docker *docker) banner(operateName string, msg string) {
	util.PrintDirectBanner([]string{constant.Docker, operateName}, msg)
}
func (docker *docker) returnBanner(msg string) util.Banner {
	return util.Banner{
		Symbols: nil,
		Msg:     msg,
	}
}

func (docker *docker) shell(cmd string, msg string) util.Shell {
	return util.Shell{
		Cmd:        cmd,
		ServerList: docker.serverList,
		Banner:     docker.returnBanner(msg),
	}
}
