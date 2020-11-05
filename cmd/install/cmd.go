package install

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/sys"
	"github.com/weiliang-ms/easyctl/util"
)

var installValidArgs = []string{"docker", "nginx", "redis"}
var (
	redisDeployMode      string
	serverListFile       string
	RedisOffline         bool
	dockerOffline        bool
	redisPort            string
	redisBindIP          string
	redisPassword        string
	redisDataDir         string
	redisLogDir          string
	redisConfigDir       string
	redisBinaryPath      string
	redisOverrideInstall bool
	redisListenAddress   string
	sourceFilePath       string
	netConnectErr        error
	installErr           error
)

func init() {
	netConnectErr = errors.New("网络连接异常,请选择离线方式安装...")
	installErr = errors.New("程序安装失败...")
	redisCmd.Flags().StringVarP(&redisPort, "port", "p", "6379", "Redis listen port 监听端口")
	redisCmd.Flags().StringVarP(&redisBindIP, "bind", "b", "0.0.0.0", "Redis bind address 监听地址")
	redisCmd.Flags().StringVarP(&redisPassword, "password", "a", "redis", "Redis password 密码")
	redisCmd.Flags().StringVarP(&redisDataDir, "data-dir", "d", "/redis/lib", "Redis persistent directory 持久化目录")
	redisCmd.Flags().StringVarP(&redisLogDir, "log-file", "", "/redis/log", "Redis logfile directory 日志目录")
	redisCmd.Flags().StringVarP(&sourceFilePath, "file", "f", "", "redis-x-x-x.tar.gz 安装包路径")
	redisCmd.Flags().StringVarP(&redisBinaryPath, "binary-path", "", "/usr/bin", "redis-* binary file path 二进制文件路径")
	redisCmd.Flags().StringVarP(&redisConfigDir, "config-file-dir", "", "/etc", "Redis 配置文件目录")
	redisCmd.Flags().BoolVarP(&redisOverrideInstall, "override", "", true, "是否覆盖安装，默认覆盖")

	// 部署模式
	redisCmd.Flags().StringVarP(&redisDeployMode, "mode", "", "single", "redis部署模式")
	// 单节点
	redisCmd.Flags().StringVarP(&redisListenAddress, "listen-address", "", "0.0.0.0", "redis监听地址")

	// 集群参数
	redisCmd.Flags().StringVarP(&serverListFile, "server-list", "", "", "ssh server连接信息配置文件路径")
	redisCmd.Flags().BoolVarP(&RedisOffline, "offline", "o", false, "offline mode")

	redisCmd.MarkFlagRequired(redisDeployMode)

	RootCmd.AddCommand(dockerCmd)
	RootCmd.AddCommand(installNginxCmd)
	RootCmd.AddCommand(redisCmd)
}

// install命令
var RootCmd = &cobra.Command{
	Use:   "install [OPTIONS] [flags]",
	Short: "install soft through easyctl",
	Example: "\neasyctl install docker" +
		"\neasyctl install nginx",
	Args: cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return installValidArgs, cobra.ShellCompDirectiveNoFileComp
	},
}

// install docker命令
var dockerCmd = &cobra.Command{
	Use:   "docker [flags]",
	Short: "install docker through easyctl",
	Example: "\neasyctl install docker 在线安装docker" +
		"\neasyctl install docker --offline --file=./docker-19.03.9.tgz 离线安装docker",
	Run: func(cmd *cobra.Command, args []string) {
		if !dockerOffline {
			installDockerOnline()
		} else {
			installDockerOffline()
		}
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

// install redis命令
var redisCmd = &cobra.Command{
	Use:   "redis [flags]",
	Short: "install redis through easyctl",
	Example: "\neasyctl install redis 在线安装redis" +
		"\neasyctl install redis --offline=true --file=./redis-5.0.5.tar.gz 离线安装redis",
	Run: func(cmd *cobra.Command, args []string) {
		var redis redis
		redis.install()
	},
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

// 安装nginx
func installNginx() {

	sys.SetNginxMirror()
	cmd := "yum -y install nginx"

	if !sys.AccessAliMirrors() {
		panic(netConnectErr)
	}
	util.ExecuteCmdAcceptResult(cmd)

	sys.CloseSeLinux(true)
	fmt.Println("[nginx]启动nginx...")
	startRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartNginx)
	fmt.Println("[nginx]设置nginx开机自启动...")
	enableRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartNginxForever)
	if startRe == nil && enableRe == nil {
		util.PrintSuccessfulMsg("nginx安装成功...")
	}

}
