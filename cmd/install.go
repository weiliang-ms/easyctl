package cmd

import (
	"easyctl/sys"
	"easyctl/util"
	"fmt"
	"github.com/spf13/cobra"
)

var installValidArgs = []string{"docker", "nginx"}

func init() {
	installCmd.AddCommand(installDockerCmd)
	installCmd.AddCommand(installNginxCmd)
	rootCmd.AddCommand(installCmd)
}

// install命令
var installCmd = &cobra.Command{
	Use:   "install [OPTIONS] [flags]",
	Short: "install some soft through easyctl",
	Example: "\neasyctl install docker" +
		"\neasyctl install nginx",
	Run: func(cmd *cobra.Command, args []string) {
	},
	ValidArgs: installValidArgs,
	Args:      cobra.ExactValidArgs(1),
}

// install docker命令
var installDockerCmd = &cobra.Command{
	Use:   "docker [flags]",
	Short: "install docker through easyctl",
	Example: "\neasyctl install docker 在线安装docker" +
		"\neasyctl install docker --offline=true --file=./v19.03.13.tar.gz 离线安装docker",
	Run: func(cmd *cobra.Command, args []string) {
		installDocker()
	},
}

// install nginx命令
var installNginxCmd = &cobra.Command{
	Use:   "nginx [flags]",
	Short: "install nginx through easyctl",
	Example: "\neasyctl install nginx 在线安装nginx" +
		"\neasyctl install nginx --offline=true --file=./nginx-1.16.0.tar.gz 离线安装nginx",
	Run: func(cmd *cobra.Command, args []string) {
		installNginx()
	},
}

// 安装docker
func installDocker() {
	fmt.Println("检测内核...")
	sys.SetAliYUM()
	cmd := "yum -y install yum-utils device-mapper-persistent-data lvm2;" +
		"yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo;" +
		"yum makecache fast;" +
		"yum -y install docker-ce"
	if sys.AccessAliMirrors() {
		_, re := util.ExecuteCmdAcceptResult(cmd)
		if re != nil {
			fmt.Println(re.Error())
		}
		sys.CloseSeLinux(true)
		fmt.Println("[docker]启动docker...")
		startRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartDocker)
		fmt.Println("[docker]设置docker开机自启动...")
		enableRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartDockerForever)
		if startRe == nil && enableRe == nil {
			util.PrintSuccessfulMsg("docker安装成功...")
		}
	}
}

// 安装nginx
func installNginx() {
	sys.SetNginxMirror()
	cmd := "yum -y install nginx"
	if sys.AccessAliMirrors() {
		_, re := util.ExecuteCmdAcceptResult(cmd)
		if re != nil {
			fmt.Println(re.Error())
		}
		sys.CloseSeLinux(true)
		fmt.Println("[nginx]启动nginx...")
		startRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartNginx)
		fmt.Println("[nginx]设置nginx开机自启动...")
		enableRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartNginxForever)
		if startRe == nil && enableRe == nil {
			util.PrintSuccessfulMsg("nginx安装成功...")
		}
	}
}
