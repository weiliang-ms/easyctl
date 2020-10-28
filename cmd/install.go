package cmd

import (
	"easyctl/constant"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/sys"
	"github.com/weiliang-ms/easyctl/util"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const RedisConfigPath = "/etc/redis.conf"

var installValidArgs = []string{"docker", "nginx", "redis"}
var (
	redisClusterMode bool
	//redisClusterNodesNum int
	nodesSSHInfoFilePath string
	offline              bool
	redisPort            string
	redisBindIP          string
	redisPassword        string
	redisDataDir         string
	redisLogDir          string
	redisConfigDir       string
	redisBinaryPath      string
	redisOverrideInstall bool
	sourceFilePath       string
	netConnectErr        error
	installErr           error
)

type redis struct {
	clusterNodes  []util.SSHInstance
	clusterEnable bool
}

func init() {
	netConnectErr = errors.New("网络连接异常,请选择离线方式安装...")
	installErr = errors.New("程序安装失败...")
	installRedisCmd.Flags().BoolVarP(&offline, "offline", "", false, "offline mode 离线模式")
	installRedisCmd.Flags().StringVarP(&redisPort, "port", "p", "6379", "Redis listen port 监听端口")
	installRedisCmd.Flags().StringVarP(&redisBindIP, "bind", "b", "0.0.0.0", "Redis bind address 监听地址")
	installRedisCmd.Flags().StringVarP(&redisPassword, "password", "a", "redis", "Redis password 密码")
	installRedisCmd.Flags().StringVarP(&redisDataDir, "data-dir", "d", "/redis/lib", "Redis persistent directory 持久化目录")
	installRedisCmd.Flags().StringVarP(&redisLogDir, "log-file", "", "/redis/log", "Redis logfile directory 日志目录")
	installRedisCmd.Flags().StringVarP(&sourceFilePath, "file", "f", "", "redis-x-x-x.tar.gz 安装包路径")
	installRedisCmd.Flags().StringVarP(&redisBinaryPath, "binary-path", "", "/usr/bin", "redis-* binary file path 二进制文件路径")
	installRedisCmd.Flags().StringVarP(&redisConfigDir, "config-file-dir", "", "/etc", "Redis 配置文件目录")

	installRedisCmd.Flags().BoolVarP(&redisOverrideInstall, "override", "", true, "是否覆盖安装，默认覆盖")

	// 集群参数
	installRedisCmd.Flags().BoolVarP(&redisClusterMode, "cluster-mode", "", true, "redis 集群模式")
	//installRedisCmd.Flags().IntVarP(&redisClusterNodesNum, "clustr-nodes-num", "", 1, "redis集群节点数量")
	installRedisCmd.Flags().StringVarP(&nodesSSHInfoFilePath, "ssh-info-file", "", "", "节点ssh连接信息配置文件路径")

	installDockerCmd.Flags().StringVarP(&sourceFilePath, "file", "f", "", "redis-x-x-x.tar.gz path")
	installDockerCmd.Flags().BoolVarP(&offline, "offline", "o", false, "offline mode")

	installCmd.AddCommand(installDockerCmd)
	installCmd.AddCommand(installNginxCmd)
	installCmd.AddCommand(installRedisCmd)
	rootCmd.AddCommand(installCmd)
}

// install命令
var installCmd = &cobra.Command{
	Use:   "install [OPTIONS] [flags]",
	Short: "install soft through easyctl",
	Example: "\neasyctl install docker" +
		"\neasyctl install nginx",
	RunE: func(cmd *cobra.Command, args []string) error {
		return parseCommand(cmd, args, installValidArgs)
	},
	Args: cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return installValidArgs, cobra.ShellCompDirectiveNoFileComp
	},
}

// install docker命令
var installDockerCmd = &cobra.Command{
	Use:   "docker [flags]",
	Short: "install docker through easyctl",
	Example: "\neasyctl install docker 在线安装docker" +
		"\neasyctl install docker --offline --file=./docker-19.03.9.tgz 离线安装docker",
	Run: func(cmd *cobra.Command, args []string) {
		if !offline {
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
var installRedisCmd = &cobra.Command{
	Use:   "redis [flags]",
	Short: "install redis through easyctl",
	Example: "\neasyctl install redis 在线安装redis" +
		"\neasyctl install redis --offline=true --file=./redis-5.0.5.tar.gz 离线安装redis",
	Run: func(cmd *cobra.Command, args []string) {
		var redis redis
		if !offline {
			installRedisOnline()
		} else {
			redis.installOffline()
		}
	},
}

// 在线安装docker
func installDockerOnline() {
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

// 离线安装docker
func installDockerOffline() {

	fmt.Println("离线安装docker...")
	docker := "tar zxvf docker-*.tgz;mv docker/* /usr/bin/"
	util.ExecuteCmdAcceptResult(docker)

	// 配置系统服务
	fmt.Println("[redis]配置redis系统服务...")
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

// 安装nginx
func installNginx() {

	sys.SetNginxMirror()
	cmd := "yum -y install nginx"

	if !sys.AccessAliMirrors() {
		panic(netConnectErr)
	}
	_, re := util.ExecuteCmdAcceptResult(cmd)
	if re != nil {
		log.Fatal(installErr)
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
}

// 离线安装redis
func (redis *redis) installOffline() {
	if redisClusterMode {
		redis.installClusterOffline()
	}
}
func (redis *redis) installClusterOffline() {
	// 解析集群节点ssh连接信息
	redis.cluterNodesInfo()

	// 检测redis安装环境
	redis.installEnvDetection()

	// 覆盖安装时执行
	redis.removeClusterConfig()

	// 拷贝源码包
	redis.copySourceCode()

	// 远程编译redis
	redis.compile()

	// 生成集群配置文件
	redis.writeClusterConfigFile()

	// 初始化运行时环境
	redis.initializeRuntimeEnv()

	// 调优
	redis.optimize()

	// 启动
	redis.startCluster()

	// 初始化集群
	redis.initializeCluster()

	// 开放端口

}

// 解析redis集群节点信息
func (redis *redis) cluterNodesInfo() {
	if _, err := os.Stat(nodesSSHInfoFilePath); err == nil {
		redis.clusterNodes = util.ReadSSHInfoFromFile(nodesSSHInfoFilePath)
	}
}

// 初始化集群
func (redis *redis) initializeCluster() {
	fmt.Println("[init] 初始化redis集群...")
	util.ExecuteCmdAcceptResult(redis.initializeClusterCmd())
}

func (redis *redis) initializeClusterCmd() (cmd string) {
	if len(redis.clusterNodes) < 2 {
		cmd = fmt.Sprintf("echo \"yes\" | redis-cli --cluster create 0.0.0.0:26379 0.0.0.0:26380 "+
			"0.0.0.0:26381 0.0.0.0:26382 0.0.0.0:26383 0.0.0.0:26384 --cluster-replicas 1 -a %s", redisPassword)
	}

	return
}

// 启动集群
func (redis *redis) startCluster() {
	fmt.Println("[start] 启动redis集群...")
	util.ExecuteCmdAcceptResult(redis.startClusterCmd())
}

func (redis *redis) startClusterCmd() (cmd string) {
	ports := []string{"26379", "26380", "26381", "26382", "26383", "26384"}
	var i int
	if len(redis.clusterNodes) < 2 {
		i = 6
	}

	for j := 0; j < i; j++ {
		cmd += fmt.Sprintf("redis-server %s/%s.conf;", redisConfigDir, ports[j])
	}

	return cmd
}

// 拷贝redis源代码
func (redis *redis) copySourceCode() {
	if len(redis.clusterNodes) == 0 {
		fmt.Println("[result] 检测到clustr-nodes-num参数为1，本地安装redis...")
		time.Sleep(3 * time.Second)
		return
	}
	// 拷贝逻辑
}

// 远程编译redis
func (redis *redis) remoteCompile(compileCmd string) {

	wg := sync.WaitGroup{}
	wg.Add(len(redis.clusterNodes))
	for i := 0; i < len(redis.clusterNodes); i++ {
		instance := redis.clusterNodes[i]
		go instance.ExecuteOriginCmdParallel(compileCmd, &wg)
	}

	wg.Wait()
}

func (redis *redis) writeClusterConfigFile() {

	fmt.Println("[init] 初始化redis集群配置文件...")
	if len(redis.clusterNodes) == 0 {
		ports := []string{"26379", "26380", "26381", "26382", "26383", "26384"}
		for _, v := range ports {
			util.OverwriteContent(fmt.Sprintf("%s/%s.conf", redisConfigDir, v), redis.config(v))
		}
	}

	if len(redis.clusterNodes) == 3 {
		ports := []string{"26379", "26380"}
		for _, v := range redis.clusterNodes {
			for _, p := range ports {
				configDir := fmt.Sprintf("%s/%s.conf", redisConfigDir, p)
				util.OriginWriteFile(configDir, []byte(constant.RedisConfigContent), v)
			}

		}
	}
}

func (redis *redis) compile() {
	fmt.Println("[compile] 开始编译redis...")
	directoryName := strings.Trim(strings.Trim(sourceFilePath, ".tar.gz"), "/")
	// sed -i "s#/usr/local#$HOME#g" src/Makefile
	compileCmd := fmt.Sprintf("tar zxvf %s && cd %s && sed -i \"s#/usr/local#%s#g\" src/Makefile && make -j $(nproc) && make install && cd ~",
		sourceFilePath, directoryName, redisBinaryPath)

	if len(redis.clusterNodes) > 0 {
		redis.remoteCompile(compileCmd)
	} else {
		util.ExecuteCmd(compileCmd)
	}

}
func (redis *redis) initializeRuntimeEnv() {
	util.ExecuteCmdIgnoreErr(redis.initEnvCmd())
}
func (redis *redis) initEnvCmd() string {
	return fmt.Sprintf("%s;mkdir --mode=0666 -p %s %s %s;",
		constant.CreateNologinUserCmd("redis"),
		redisDataDir, redisLogDir, redisConfigDir)
}

// 配置redis
func (redis *redis) config(port string) string {
	// todo 校验参数合法性
	content := strings.ReplaceAll(constant.RedisConfigContent, "26379", port)
	strings.ReplaceAll(content, "requirepass redis", fmt.Sprintf("requirepass %s", redisPassword))
	strings.ReplaceAll(content, "dir /redis/lib", fmt.Sprintf("dir %s", redisDataDir))
	strings.ReplaceAll(content, "logfile /redis/log", fmt.Sprintf("dir %s", redisLogDir))

	//fmt.Println("[redis]配置overcommit_memory...")
	//
	return content
}

func (redis *redis) optimize() {
	util.ExecuteCmd("echo \"vm.overcommit_memory = 1\" >> /etc/sysctl.conf;sysctl -p")
}

func (redis *redis) start() {
	fmt.Println("[redis] 启动redis...")
	startRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartRedis)
	fmt.Println("[redis] 配置redis为系统服务...")
	enableRe, _ := util.ExecuteCmd(sys.SystemInfoObject.ServiceAction.StartRedisForever)
	if startRe == nil && enableRe == nil {
		util.PrintSuccessfulMsg("[success] redis安装成功...")
	}
}

func (redis *redis) installEnvDetection() {
	fmt.Println("[check] 检测redis编译环境...")
	gcc := true
	re := ""
	for _, v := range redis.clusterNodes {
		if !util.RemoteHostYumDetection(v) {
			re += fmt.Sprintf("节点:%s yum检测失败，请配置好yum源\n", v.Host)
			gcc = false
		}
	}
	if !gcc {
		log.Fatal("[failed] redis安装环境检测失败...")
	}
}

func (redis *redis) openFirewallPort() {

}

func (redis *redis) openFirewallPortCmd() {
	ports := []int{}
	var i int
	if len(redis.clusterNodes) == 0 || len(redis.clusterNodes) == 1 {
		i = 6
	} else if len(redis.clusterNodes) == 3 {
		i = 2
	} else if len(redis.clusterNodes) == 2 {
		i = 3
	} else {
		fmt.Sprintf("节点数量参数与不匹配")
	}

	for j := 0; j < i; j++ {
		ports = append(ports, 26379+j)
		ports = append(ports, 26379+j+10000)
	}
}

// 删除集群配置
func (redis *redis) removeClusterConfig() {
	if redisOverrideInstall {
		util.ExecuteCmdIgnoreErr(redis.removeClusterConfigCmd())
	}
}

func (redis *redis) removeClusterConfigCmd() string {
	return fmt.Sprintf("pkill -9 redis;rm -f %s/nodes-263* %s/* %s/*",
		redisConfigDir, redisDataDir, redisLogDir)
}
