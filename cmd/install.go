package cmd

import (
	"easyctl/sys"
	"easyctl/util"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

const RedisConfigPath = "/etc/redis.conf"

var installValidArgs = []string{"docker", "nginx", "redis"}
var (
	Offline       bool
	RedisPort     string
	RedisBindIP   string
	RedisPassword string
	RedisDataDir  string
	RedisLogDir   string
	netConnectErr error
	installErr    error
)

func init() {
	netConnectErr = errors.New("网络连接异常,请选择离线方式安装...")
	installErr = errors.New("程序安装失败...")
	installRedisCmd.Flags().BoolVarP(&Offline, "offline", "o", false, "offline mode")
	installRedisCmd.Flags().StringVarP(&RedisPort, "port", "p", "6379", "Redis listen port")
	installRedisCmd.Flags().StringVarP(&RedisBindIP, "bind", "b", "0.0.0.0", "Redis bind address")
	installRedisCmd.Flags().StringVarP(&RedisPassword, "password", "a", "redis", "Redis password")
	installRedisCmd.Flags().StringVarP(&RedisDataDir, "data", "d", "/var/lib/redis", "Redis persistent directory")
	installRedisCmd.Flags().StringVarP(&RedisDataDir, "log-file", "l", "/var/log/redis", "Redis logfile directory")
	installCmd.AddCommand(installDockerCmd)
	installCmd.AddCommand(installNginxCmd)
	installCmd.AddCommand(installRedisCmd)
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

// install redis命令
var installRedisCmd = &cobra.Command{
	Use:   "redis [flags]",
	Short: "install redis through easyctl",
	Example: "\neasyctl install redis 在线安装redis" +
		"\neasyctl install redis --offline=true --file=./redis-5.0.5.tar.gz 离线安装redis",
	Run: func(cmd *cobra.Command, args []string) {
		if !Offline {
			installRedisOnline()
		}
	},
}

// 安装docker
func installDocker() {
	fmt.Println("检测内核...")
	if !sys.AccessAliMirrors() {
		panic(netConnectErr)
	}

	sys.SetAliYUM()
	install := "yum -y install yum-utils device-mapper-persistent-data lvm2;" +
		"yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo;" +
		"yum makecache fast;" +
		"yum -y install docker-ce"

	//
	_, re := util.ExecuteCmdAcceptResult(install)
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

// 安装nginx
func installNginx() {

	sys.SetNginxMirror()
	cmd := "yum -y install nginx"

	if !sys.AccessAliMirrors() {
		panic(netConnectErr)
	}
	_, re := util.ExecuteCmdAcceptResult(cmd)
	if re != nil {
		panic(installErr)
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

// 在线安装redis
func installRedisOnline() {

	installCmd := "yum -y install redis"
	if !sys.AccessAliMirrors() {
		panic(netConnectErr)
	}

	sys.SetAliYUM()
	_, re := util.ExecuteCmdAcceptResult(installCmd)
	if re != nil {
		panic(installErr)
	}

	configRedis()
	fmt.Println("[redis]启动redis...")
	startRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartRedis)
	fmt.Println("[redis]配置redis为系统服务...")
	enableRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartRedisForever)
	if startRe == nil && enableRe == nil {
		util.PrintSuccessfulMsg("redis安装成功...")
	}
}

// 配置redis
func configRedis() {
	modifyBgsaveErr := fmt.Sprintf("sed -i \"s#stop-writes-on-bgsave-error yes#stop-writes-on-bgsave-error no#g\" %s", RedisConfigPath)
	modifyDaemonCmd := fmt.Sprintf("sed -i 's#daemonize no#daemonize yes#g' %s", RedisConfigPath)
	modifyBindCmd := fmt.Sprintf("sed -i \"s#port 6379#port %s#g\" %s", RedisPort, RedisConfigPath)
	modifyListenPortCmd := fmt.Sprintf("sed -i \"s#bind 127.0.0.1#bind %s#g\" %s", RedisBindIP, RedisConfigPath)
	modifyPassword := fmt.Sprintf("echo \"requirepass %s\" >> %s", RedisPassword, RedisConfigPath)
	modifyDataDir := fmt.Sprintf("sed -i \"s#dir ./#dir %s#g\" %s", RedisDataDir, RedisConfigPath)
	modifyLogDir := fmt.Sprintf("sed -i \"s#logfile \\\"\\\"#logfile %s/redis.log#g\" %s", RedisLogDir, RedisConfigPath)
	fmt.Println("[redis]配置redis...")
	util.ExecuteCmd(modifyBgsaveErr)
	util.ExecuteCmd(modifyDaemonCmd)
	util.ExecuteCmd(modifyBindCmd)
	util.ExecuteCmd(modifyListenPortCmd)
	util.ExecuteCmd(modifyPassword)
	util.ExecuteCmd(modifyDataDir)
	util.ExecuteCmd(modifyLogDir)
	fmt.Println("[redis]创建redis工作目录...")
	util.ExecuteCmd(fmt.Sprintf("mkdir -p %s;mkdir -p %s", RedisLogDir, RedisDataDir))
	fmt.Println("[redis]配置overcommit_memory...")
	util.ExecuteCmd("echo \"vm.overcommit_memory = 1\" >> /etc/sysctl.conf;sysctl -p")
}
