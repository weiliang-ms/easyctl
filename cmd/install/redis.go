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
	originDeploy bool
	offline      bool
	mode         mode
	bind         string
}

const (
	single                     = "single"
	clusterOneNode             = "clusterOneNode"
	startRedisCluster          = "启动redis集群"
	initRedisCluster           = "初始化redis集群"
	initRedisRuntimeEnv        = "初始化redis运行时环境"
	modifyFirewallRule         = "调整防火墙策略"
	redisOptimize              = "redis调优"
	removeOldConfig            = "删除redis历史版本文件"
	compileRedisEnvDetection   = "检测redis编译环境"
	dependenceDetectionNotPass = "依赖检测未通过,请检查yum配置"
	generateRedisConfigFile    = "生成redis配置文件"
	startToCompileRedis        = "开始编译redis"
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
	if redisClusterMode {
		redis.installClusterOffline()
	}
}
func (redis *redis) installClusterOffline() {

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

	// 远程编译redis
	redis.compile()

	// 生成集群配置文件
	redis.writeClusterConfigFile()

	// 初始化运行时环境
	redis.initializeRuntimeEnv()

	// 调优
	redis.optimize()

	// 启动
	redis.start()

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
	if redisClusterMode && len(redis.clusterNodes) == 0 && redisBindIP == "" {
		log.Fatal("当前模式必须指定--bind参数...")
	}
}

func (redis *redis) deployMode() *redis {
	if len(redis.clusterNodes) > 0 {
		redis.originDeploy = true
	}

	return redis
}

func (redis *redis) nodeMode() *redis {
	// 解析是否单点
	if redisSingleMode {
		redis.mode = single
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
	if _, err := os.Stat(nodesSSHInfoFilePath); err == nil {
		redis.clusterNodes = util.ReadSSHInfoFromFile(nodesSSHInfoFilePath)
	}
}

// 初始化集群
func (redis *redis) initializeCluster() {
	redis.shell(initRedisCluster, redis.initializeClusterCmd())
}

func (redis *redis) initializeClusterCmd() (cmd string) {
	if len(redis.clusterNodes) < 2 {
		if len(redis.clusterNodes) == 1 {
			redisBindIP = redis.clusterNodes[0].Host
		}
		cmd = fmt.Sprintf("echo \"yes\" | redis-cli --cluster create %s:26379 %s:26380 %s:26381 %s:26382 %s:26383 %s:26384 --cluster-replicas 1 -a %s", redisBindIP,
			redisBindIP, redisBindIP, redisBindIP, redisBindIP, redisBindIP, redisPassword)
	} else if len(redis.clusterNodes) == 2 {

	}

	return
}

// 启动集群
func (redis *redis) start() {
	switch redis.mode {
	case single:
		redis.startSingle()
	case clusterOneNode:
		redis.startClusterOneNode()
	default:
		log.Fatal("暂不支持这种部署方式...")
	}
}

// 单节点非集群实例
func (redis *redis) startSingle() {
	fmt.Println("启动单节点实例...")
}

// 单节点集群实例
func (redis *redis) startClusterOneNode() {
	redis.shell(startRedisCluster, redis.startClusterCmd())
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

func (redis *redis) banner(msg string, address string) {
	var location string

	if redis.originDeploy {
		location = constant.Origin
	} else {
		location = constant.Local
		address = constant.LoopbackAddress
	}

	fmt.Printf("%s %s...\n",
		util.PrintOrangeMulti([]string{constant.Redis, location, address}), msg)
}

// 配置systemd service
func (redis *redis) service() {

}

func (redis *redis) serviceCmd() (cmd string) {
	if redis.mode == single {
		cmd = fmt.Sprintf("chmod +x /etc/rc.local;ls %s/263*.conf|xargs -n1 install-server", redisConfigDir)
	}

	return
}

// 拷贝redis源代码
func (redis *redis) scpSourceCode() {

	if !redis.originDeploy {
		time.Sleep(3 * time.Second)
		return
	} else {
		file, _ := os.OpenFile(sourceFilePath, os.O_RDONLY, 0666)
		b, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err.Error())
		}
		for _, v := range redis.clusterNodes {
			fmt.Printf("%s 拷贝源码包至%s:~/%s...\n", util.PrintOrangeMulti([]string{constant.Redis, constant.Origin, v.Host}), v.Host, file.Name())
			util.OriginWriteFile(fmt.Sprintf("%s/%s", util.HomeDir(v), file.Name()), b, v)
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
		fmt.Printf("%s 开始编译redis...\n", util.PrintOrangeMulti([]string{constant.Redis, constant.Origin, instance.Host}))
		go instance.ExecuteOriginCmdParallel(compileCmd, &wg)
	}

	wg.Wait()
}

func (redis *redis) writeClusterConfigFile() {

	if len(redis.clusterNodes) == 0 {
		ports := []string{"26379", "26380", "26381", "26382", "26383", "26384"}
		for _, v := range ports {
			redis.banner(fmt.Sprintf("%s: %s/%s.conf...", generateRedisConfigFile, redisConfigDir, v), constant.LoopbackAddress)
			util.OverwriteContent(fmt.Sprintf("%s/%s.conf", redisConfigDir, v), redis.config(v))
		}
	}
	if len(redis.clusterNodes) == 1 {
		ports := []string{"26379", "26380", "26381", "26382", "26383", "26384"}
		for _, v := range ports {
			configFilePath := fmt.Sprintf("%s/%s.conf", redisConfigDir, v)
			content := strings.ReplaceAll(constant.RedisConfigContent, "26379", v)
			fmt.Printf("%s 初始化redis集群配置文件: %s...\n", util.PrintOrangeMulti([]string{constant.Redis, constant.Origin, redis.clusterNodes[0].Host}), configFilePath)
			util.OriginWriteFile(configFilePath, []byte(content), redis.clusterNodes[0])
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

	if !redis.originDeploy {
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
						util.PrintRedMulti([]string{constant.Error, constant.Origin, v.Host}),
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
	num := len(redis.clusterNodes)
	if num == 0 || num == 1 {
		i = 6
	} else if num == 3 {
		i = 2
	} else if num == 2 {
		i = 3
	} else {
		fmt.Sprintf("节点数量参数与不匹配")
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

	if !redis.originDeploy {
		redis.banner(msg, constant.LoopbackAddress)
		util.ExecuteCmdAcceptResult(cmd)
	} else {
		for _, v := range redis.clusterNodes {
			redis.banner(msg, v.Host)
			v.ExecuteOriginCmd(cmd)
		}
	}
}
