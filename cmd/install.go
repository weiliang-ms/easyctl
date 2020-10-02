package cmd

import (
	"easyctl/sys"
	"easyctl/util"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const RedisConfigPath = "/etc/redis.conf"

var installValidArgs = []string{"docker", "nginx", "redis"}
var (
	offline         bool
	redisPort       string
	redisBindIP     string
	redisPassword   string
	redisDataDir    string
	redisLogDir     string
	redisBinaryPath string
	filePath        string
	netConnectErr   error
	installErr      error
)

func init() {
	netConnectErr = errors.New("网络连接异常,请选择离线方式安装...")
	installErr = errors.New("程序安装失败...")
	installRedisCmd.Flags().BoolVarP(&offline, "offline", "o", false, "offline mode")
	installRedisCmd.Flags().StringVarP(&redisPort, "port", "p", "6379", "Redis listen port")
	installRedisCmd.Flags().StringVarP(&redisBindIP, "bind", "b", "0.0.0.0", "Redis bind address")
	installRedisCmd.Flags().StringVarP(&redisPassword, "password", "a", "redis", "Redis password")
	installRedisCmd.Flags().StringVarP(&redisDataDir, "data", "d", "/var/lib/redis", "Redis persistent directory")
	installRedisCmd.Flags().StringVarP(&redisLogDir, "log-file", "", "/var/log/redis", "Redis logfile directory")
	installRedisCmd.Flags().StringVarP(&filePath, "file", "f", "", "redis-x-x-x.tar.gz path")
	installRedisCmd.Flags().StringVarP(&redisBinaryPath, "binary-path", "", "/usr/bin/", "redis-* binary file path")

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
		"\neasyctl install docker --offline --file=./v19.03.13.tar.gz 离线安装docker",
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
		if !offline {
			installRedisOnline()
		} else {
			installRedisOffline()
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
	startRedis()
}

// 在线安装redis
func installRedisOffline() {

	gcc := "rpm -qa|grep \"^gcc\";echo $?"
	gccRe, _ := util.ExecuteCmdAcceptResult(gcc)
	if gccRe != "0" {
		installGcc, _ := util.ExecuteCmdAcceptResult("yum install -y gcc;echo $?")
		if installGcc != "0" {
			panic(errors.New("安装gcc失败，请配置好yum源..."))
		}
	}

	_, err := os.Stat(filePath)
	if err != nil {
		panic(errors.New("访问redis源码包失败"))
	}

	compileRedis()
	startRedis()
}

func compileRedis() {
	fmt.Println("[redis]编译redis...")
	directoryName := strings.Trim(strings.Trim(filePath, ".tar.gz"), "/")
	util.ExecuteCmdAcceptResult("tar zxvf " + filePath)
	util.ExecuteCmdAcceptResult(fmt.Sprintf("cd %s;make;make install;cp redis.conf %s;cd -", directoryName, RedisConfigPath))
	util.ExecuteCmd(fmt.Sprintf("cp /usr/local/bin/redis-* %s", redisBinaryPath))

	// 创建redis用户
	fmt.Println("[redis]创建redis用户...")
	sys.AddUser("redis", "", false)

	fmt.Println("[redis]配置redis...")
	configRedis()
	// 配置系统服务
	fmt.Println("[redis]配置redis系统服务...")
	sys.ConfigService("redis")

}

// 配置redis
func configRedis() {
	// todo 校验参数合法性
	modifyBgsaveErr := fmt.Sprintf("sed -i \"s#stop-writes-on-bgsave-error yes#stop-writes-on-bgsave-error no#g\" %s", RedisConfigPath)
	modifyDaemonCmd := fmt.Sprintf("sed -i 's#daemonize no#daemonize yes#g' %s", RedisConfigPath)
	modifyBindCmd := fmt.Sprintf("sed -i \"s#port 6379#port %s#g\" %s", redisPort, RedisConfigPath)
	modifyListenPortCmd := fmt.Sprintf("sed -i \"s#bind 127.0.0.1#bind %s#g\" %s", redisBindIP, RedisConfigPath)
	modifyPassword := fmt.Sprintf("echo \"requirepass %s\" >> %s", redisPassword, RedisConfigPath)
	modifyDataDir := fmt.Sprintf("sed -i \"s#dir ./#dir %s#g\" %s", redisDataDir, RedisConfigPath)
	modifyLogDir := fmt.Sprintf("sed -i \"s#logfile \\\"\\\"#logfile %s/redis.log#g\" %s", redisLogDir, RedisConfigPath)
	util.ExecuteCmd(modifyBgsaveErr)
	util.ExecuteCmd(modifyDaemonCmd)
	util.ExecuteCmd(modifyBindCmd)
	util.ExecuteCmd(modifyListenPortCmd)
	util.ExecuteCmd(modifyPassword)
	util.ExecuteCmd(modifyDataDir)
	util.ExecuteCmd(modifyLogDir)
	fmt.Println("[redis]创建redis工作目录...")
	util.ExecuteCmd(fmt.Sprintf("mkdir -p %s;mkdir -p %s", redisLogDir, redisDataDir))
	fmt.Println("[redis]配置overcommit_memory...")
	util.ExecuteCmd("echo \"vm.overcommit_memory = 1\" >> /etc/sysctl.conf;sysctl -p")
}

func startRedis() {
	fmt.Println("[redis]启动redis...")
	startRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartRedis)
	fmt.Println("[redis]配置redis为系统服务...")
	enableRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartRedisForever)
	if startRe == nil && enableRe == nil {
		util.PrintSuccessfulMsg("redis安装成功...")
	}
}
