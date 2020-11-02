package install

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/constant"
	"github.com/weiliang-ms/easyctl/util"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type mode string

type redis struct {
	clusterNodes []util.SSHInstance
	cluster      bool
	remoteDeploy bool
	offline      bool
	mode         mode
	bind         string
}

const (
	singleLocal                = "singleLocal"
	singleRemote               = "singleRemote"
	clusterOneNode             = "clusterOneNode"
	startRedis                 = "启动redis服务"
	initRedisCluster           = "初始化redis集群"
	initRedisRuntimeEnv        = "初始化redis运行时环境"
	modifyFirewallRule         = "调整防火墙策略"
	redisOptimize              = "redis调优"
	removeOldConfig            = "删除redis历史版本文件"
	compileRedisEnvDetection   = "检测redis编译环境"
	dependenceDetectionNotPass = "依赖检测未通过,请检查yum配置"
	generateRedisConfigFile    = "生成redis配置文件"
	startToCompileRedis        = "开始编译redis"
	configService              = "配置开机自启动"
)

func (redis *redis) install() {
	if RedisOffline {
		redis.installOffline()
	} else {
		installRedisOnline()
	}
}

// 在线安装redis
func installRedisOnline() {

}

// 离线安装redis

func (redis *redis) installOffline() {

	// 解析集群节点ssh连接信息
	redis.cluterNodesInfo()

	// 解析参数
	redis.parseFlags()

	// 参数检测
	redis.require()

	// 检测redis安装环境
	redis.compileEnvDetection()

	// 拷贝源码包
	redis.scpSourceCode()

	// 覆盖安装时执行
	redis.removeClusterConfig()

	// 编译redis
	redis.compile()

	// 生成配置文件
	redis.generateConfigFile()

	// 初始化运行时环境
	redis.initializeRuntimeEnv()

	// 调优
	redis.optimize()

	// 启动
	redis.start()

	// 配置系统级服务
	redis.service()

	// 初始化集群
	redis.initializeCluster()

	// 开放端口
	redis.openFirewallPort()
}

func (redis *redis) parseFlags() {
	redis.nodeMode().deployMode()
}

func (redis *redis) require() {
	// todo 检测参数合法性逻辑
	if redisClusterMode && len(redis.clusterNodes) == 0 && redisListenAddress == "" {
		log.Fatal("当前模式必须指定--listen-address参数...")
	}
}

func (redis *redis) deployMode() *redis {
	if len(redis.clusterNodes) > 0 {
		redis.remoteDeploy = true
	}

	return redis
}

func (redis *redis) nodeMode() *redis {
	// 解析是否单点
	if redisSingleMode && len(redis.clusterNodes) > 0 {
		redis.mode = singleRemote
	} else if redisSingleMode && len(redis.clusterNodes) == 0 {
		redis.mode = singleLocal
	}

	// 远端单机集群
	if redisClusterMode && len(redis.clusterNodes) == 1 {
		redis.mode = clusterOneNode
	}

	// 本地单机集群
	if redisClusterMode && len(redis.clusterNodes) == 0 {
		redis.mode = clusterOneNode
	}

	return redis
}

// 解析redis集群节点信息
func (redis *redis) cluterNodesInfo() {
	fmt.Printf("%s 解析ssh配置清单...\n", util.PrintOrange(constant.Redis))
	if _, err := os.Stat(serverList); err == nil {
		redis.clusterNodes = util.ReadSSHInfoFromFile(serverList)
	}
}

// 初始化集群
func (redis *redis) initializeCluster() {
	if redisClusterMode {
		redis.shell(initRedisCluster, redis.initializeClusterCmd())
	}
}

func (redis *redis) initializeClusterCmd() (cmd string) {
	if redis.mode == clusterOneNode {
		cmd = redis.initializeOneNodeClusterCmd()
	} else if len(redis.clusterNodes) == 2 {

	}

	return
}

func (redis *redis) initializeOneNodeClusterCmd() (cmd string) {
	cmd = "echo \"yes\" | redis-cli --cluster create"
	ports := []string{"26379", "26380", "26381", "26382", "26383", "26384"}
	for _, v := range ports {
		cmd += fmt.Sprintf(" %s:%s", redisBindIP, v)
	}
	cmd += fmt.Sprintf(" --cluster-replicas 1 -a %s", redisPassword)
	return
}

// 启动redis服务
func (redis *redis) start() {
	redis.shell(startRedis, redis.startCmd())
}

func (redis *redis) startCmd() (cmd string) {
	ports := []string{"26379", "26380", "26381", "26382", "26383", "26384"}
	var i int
	if redis.mode == clusterOneNode {
		i = 6
	}

	if redis.mode == singleLocal || redis.mode == singleRemote {
		return fmt.Sprintf("%s/redis-server %s/redis.conf",
			redisBinaryPath, redisConfigDir)
	}

	for j := 0; j < i; j++ {
		cmd += fmt.Sprintf("redis-server %s/%s.conf;", redisConfigDir, ports[j])
	}

	return cmd
}

func (redis *redis) banner(msg string, address string) {
	var location string

	if redis.remoteDeploy {
		location = constant.Remote
	} else {
		location = constant.Local
		address = constant.LoopbackAddress
	}

	fmt.Printf("%s %s...\n",
		util.PrintOrangeMulti([]string{constant.Redis, location, address}), msg)
}

// 配置systemd service
func (redis *redis) service() {
	redis.shell(configService, redis.serviceCmd())
}

func (redis *redis) serviceCmd() (cmd string) {
	if redis.mode == singleLocal || redis.mode == singleRemote {
		cmd = fmt.Sprintf("%s %s;sed -i \"/redis-server/d\" %s;echo \"%s/redis-server %s/redis.conf\" >> %s",
			constant.ChmodX, constant.EtcRcLocal, constant.EtcRcLocal, redisBinaryPath, redisConfigDir, constant.EtcRcLocal)
	}
	if redis.mode == clusterOneNode {
		cmd = fmt.Sprintf("%s %s;sed -i \"/redis-server/d\" %s;echo \"%s/redis-server %s/263*.conf\" >> %s",
			constant.ChmodX, constant.EtcRcLocal, constant.EtcRcLocal, redisBinaryPath, redisConfigDir, constant.EtcRcLocal)
	}

	return
}

// 拷贝redis源代码
func (redis *redis) scpSourceCode() {

	if !redis.remoteDeploy {
		time.Sleep(3 * time.Second)
		return
	} else {
		file, _ := os.OpenFile(sourceFilePath, os.O_RDONLY, 0666)
		b, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err.Error())
		}
		for _, v := range redis.clusterNodes {
			fmt.Printf("%s 拷贝源码包至%s:~/%s...\n", util.PrintOrangeMulti([]string{constant.Redis, constant.Remote, v.Host}), v.Host, file.Name())
			util.RemoteWriteFile(fmt.Sprintf("%s/%s", util.HomeDir(v), file.Name()), b, v)
		}
	}
	// 拷贝逻辑
}

// 远程编译redis
func (redis *redis) remoteCompile(compileCmd string) {

	wg := sync.WaitGroup{}
	wg.Add(len(redis.clusterNodes))
	for i := 0; i < len(redis.clusterNodes); i++ {
		instance := redis.clusterNodes[i]
		fmt.Printf("%s 开始编译redis...\n", util.PrintOrangeMulti([]string{constant.Redis, constant.Remote, instance.Host}))
		go instance.ExecuteOriginCmdParallel(compileCmd, &wg)
	}

	wg.Wait()
}

func (redis *redis) generateConfigFile() {

	if redis.mode == singleLocal {
		redis.banner(generateRedisConfigFile, redisListenAddress)
		util.OverwriteContent(fmt.Sprintf("%s/redis.conf", redisConfigDir), redis.config("6379"))
	}

	if redis.mode == singleRemote {
		redis.banner(generateRedisConfigFile, redisListenAddress)
		util.RemoteWriteFile(fmt.Sprintf("%s/redis.conf", redisConfigDir), []byte(redis.config("6379")), redis.clusterNodes[0])
	}

	if redis.mode == clusterOneNode && !redis.remoteDeploy {
		ports := []string{"26379", "26380", "26381", "26382", "26383", "26384"}
		for _, v := range ports {
			redis.banner(fmt.Sprintf("%s: %s/%s.conf...", generateRedisConfigFile, redisConfigDir, v), constant.LoopbackAddress)
			util.OverwriteContent(fmt.Sprintf("%s/%s.conf", redisConfigDir, v), redis.config(v))
		}
	}
	if redis.mode == clusterOneNode && redis.remoteDeploy {
		ports := []string{"26379", "26380", "26381", "26382", "26383", "26384"}
		for _, v := range ports {
			configFilePath := fmt.Sprintf("%s/%s.conf", redisConfigDir, v)
			content := strings.ReplaceAll(constant.RedisConfigContent, "26379", v)
			fmt.Printf("%s 初始化redis集群配置文件: %s...\n", util.PrintOrangeMulti([]string{constant.Redis, constant.Remote, redis.clusterNodes[0].Host}), configFilePath)
			util.RemoteWriteFile(configFilePath, []byte(content), redis.clusterNodes[0])
		}
	}

	if len(redis.clusterNodes) == 3 {
		ports := []string{"26379", "26380"}
		for _, v := range redis.clusterNodes {
			for _, p := range ports {
				configDir := fmt.Sprintf("%s/%s.conf", redisConfigDir, p)
				util.RemoteWriteFile(configDir, []byte(constant.RedisConfigContent), v)
			}
		}
	}
}

// 编译
func (redis *redis) compile() {
	directoryName := util.CutCharacter(strings.TrimSuffix(sourceFilePath, ".tar.gz"), []string{"./", ""})
	compileCmd := fmt.Sprintf("tar zxvf %s && cd %s && sed -i \"s#\\$(PREFIX)/bin#%s#g\" src/Makefile && make -j $(nproc) && make install && cd ~",
		sourceFilePath, directoryName, redisBinaryPath)

	redis.shell(startToCompileRedis, compileCmd)

}
func (redis *redis) initializeRuntimeEnv() {
	redis.shell(initRedisRuntimeEnv, redis.initEnvCmd())
}
func (redis *redis) initEnvCmd() string {
	return fmt.Sprintf("%s;mkdir --mode=0666 -p %s %s %s;",
		constant.CreateNologinUserCmd("redis"),
		redisDataDir, redisLogDir, redisConfigDir)
}

// 配置redis
func (redis *redis) config(port string) (content string) {

	content = strings.ReplaceAll(constant.RedisConfigContent, "26379", port)
	strings.ReplaceAll(content, "requirepass redis", fmt.Sprintf("requirepass %s", redisPassword))
	strings.ReplaceAll(content, "dir /redis/lib", fmt.Sprintf("dir %s", redisDataDir))
	strings.ReplaceAll(content, "logfile /redis/log", fmt.Sprintf("logfile %s", redisLogDir))
	if redis.mode == singleRemote || redis.mode == singleLocal {
		fmt.Println("替换")
		content = strings.ReplaceAll(content, "cluster-enabled yes", "cluster-enabled no")
	}
	return
}

// 参数调优
func (redis *redis) optimize() {
	redis.shell(redisOptimize, redis.optimizeCmd())
}

func (redis *redis) optimizeCmd() string {
	return fmt.Sprintf("%s && %s && %s",
		constant.RootDetectionCmd, constant.OverCommitMemoryOptimizeCmd, constant.LimitOptimizeCmd)
}

func (redis *redis) compileEnvDetection() {
	search := fmt.Sprintf("%s %s", constant.RpmSearch, constant.Gcc)
	install := fmt.Sprintf("%s %s", constant.YumInstall, "gcc")

	if !redis.remoteDeploy {
		redis.banner(compileRedisEnvDetection, constant.Local)
		if !util.ExecuteIgnoreStd(search) {
			if !util.ExecuteIgnoreStd(install) {
				fmt.Printf("%s %s...\n",
					util.PrintRedMulti([]string{constant.Error, constant.Local, constant.LoopbackAddress}),
					dependenceDetectionNotPass)
			}
		}
	} else {
		for _, v := range redis.clusterNodes {
			redis.banner(compileRedisEnvDetection, v.Host)
			if !v.ExecuteOriginCmdIgnoreRe(search) {
				if !v.ExecuteOriginCmdIgnoreRe(install) {
					fmt.Printf("%s 节点：%s %s...\n",
						util.PrintRedMulti([]string{constant.Error, constant.Remote, v.Host}),
						v.Host,
						dependenceDetectionNotPass)
					os.Exit(1)
				}
			}
		}
	}
}

// 开启端口
func (redis *redis) openFirewallPort() {
	redis.shell(modifyFirewallRule, redis.firewallPortsCmd())
}

func (redis *redis) firewallPortsCmd() string {

	var i int
	var ports []int

	if redis.mode == clusterOneNode {
		i = 6
	}

	if redis.mode == singleLocal || redis.mode == singleRemote {
		ports = append(ports, 6379)
	}

	for j := 0; j < i; j++ {
		ports = append(ports, 26379+j)
		ports = append(ports, 26379+j+10000)
	}

	var cmd6, cmd7 string
	for _, v := range ports {
		cmd7 += fmt.Sprintf("\nfirewall-cmd --zone=public --add-port=%d/tcp --permanent;", v)
		cmd6 += fmt.Sprintf("\niptables -I INPUT -p tcp -m state --state NEW -m tcp --dport %d -j ACCEPT;", v)
	}

	cmd6 = fmt.Sprintf("%s && %s [ -f /etc/rc.d/init.d/iptables ] && /etc/rc.d/init.d/iptables save && service iptables restart", constant.Redhat6, cmd6)
	cmd7 = fmt.Sprintf("%s && %s firewall-cmd --reload", constant.Redhat7, cmd7)

	return fmt.Sprintf("%s && %s;\n%s", constant.RootDetectionCmd, cmd6, cmd7)
}

// 删除集群配置
func (redis *redis) removeClusterConfig() {
	redis.shell(removeOldConfig, redis.removeClusterConfigCmd())
}

func (redis *redis) removeClusterConfigCmd() string {
	return fmt.Sprintf("pkill -9 redis;rm -f %s/nodes-263* %s/* %s/*",
		redisConfigDir, redisDataDir, redisLogDir)
}

func (redis *redis) shell(msg string, cmd string) {

	if !redis.remoteDeploy {
		redis.banner(msg, constant.LoopbackAddress)
		util.ExecuteCmdAcceptResult(cmd)
	} else {
		for _, v := range redis.clusterNodes {
			redis.banner(msg, v.Host)
			v.ExecuteOriginCmd(cmd)
		}
	}
}
