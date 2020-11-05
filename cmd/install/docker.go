package install

import (
	"easyctl/constant"
	"easyctl/sys"
	"easyctl/util"
	"fmt"
	"log"
	"os"
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
	generateDockerServiceFile = "生成docker service配置文件"
	reloadUnitService         = "重载系统service unit配置文件"
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

	// 解压
	docker.decompress()

	// 写unit服务配置
	docker.writeService()

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
	util.Shell{
		Cmd:        fmt.Sprintf("%s && tar zxvf docker-*.tgz;mv docker/* /usr/bin/", constant.RootDetectionCmd),
		ServerList: docker.serverList,
		Banner:     util.Banner{Symbols: []string{constant.Docker}, Msg: decompressDockerPackage},
	}.Shell()
}

// 拷贝
func (docker *docker) scp() {
	// 拷贝docker安装包
	if len(docker.serverList) == 0 {
		time.Sleep(3 * time.Second)
		return
	} else {
		util.ScpHome(util.Banner{
			Symbols: []string{constant.Docker, constant.Scp},
			Msg:     "",
		}, dockerPackageFile, docker.serverList)
	}
}

// 配置系统服务
func (docker *docker) writeService() {
	util.WriteFile(docker.serviceFilePath, []byte(docker.serviceContent), docker.serverList)
}

func (docker *docker) reload() {
	if len(docker.serverList) == 0 {
		docker.banner(constant.Service, generateDockerServiceFile)
		util.OverwriteContent(docker.serviceFilePath, docker.serviceContent)

		docker.banner(constant.Reload, reloadUnitService)
		util.ExecuteCmdAcceptResult(constant.DeamonReload)
	} else {
		util.Shell{
			Cmd:        "",
			ServerList: nil,
			Banner:     util.Banner{},
		}.Shell()
	}
}

// 离线安装docker
func installDockerOffline() {

	// 配置系统服务
	fmt.Println("[install]配置redis系统服务...")
	sys.ConfigService("docker")

	sys.CloseSeLinux(true)
	fmt.Println("[docker]启动docker...")
	startRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartDocker)
	fmt.Println("[docker]设置docker开机自启动...")
	enableRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartDockerForever)
	if startRe == nil && enableRe == nil {
		util.PrintSuccessfulMsg("docker安装成功...")
	}
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
	util.PrintBanner([]string{constant.Docker, operateName}, msg)
}
