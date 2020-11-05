package install

import (
	"easyctl/sys"
	"easyctl/util"
	"fmt"
)

type docker struct {
	mode string
}

func (docker *docker) Install() {

}

// 离线安装docker
func installDockerOffline() {

	fmt.Println("离线安装docker...")
	docker := "tar zxvf docker-*.tgz;mv docker/* /usr/bin/"
	util.ExecuteCmdAcceptResult(docker)

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
