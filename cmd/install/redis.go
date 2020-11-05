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
	clusterNodes    []util.Server
	cluster         bool
	remoteDeploy    bool
	offline         bool
	mode            mode
	localCluster    bool
	clusterNodesNum int
	bind            string
	connectAddress  []string
	listenAddress   string
}

const (
	single                     = "single"
	cluster                    = "cluster"
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
	clusterStatus              = "查看redis集群状态"
	redisConnectionInfo        = "redis连接信息"
	parallel                   = "parallel"
	springRedisPassword        = "--spring.redis.password"
	springRedisHost            = "--spring.redis.host"
	springRedisPort            = "--spring.redis.port"
	springRedisClusterNodes    = "--spring.redis.cluster.nodes"
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

	// 状态
	redis.status()

	// 连接方式
	// 状态
	redis.connection()
}

func (redis *redis) parseFlags() {
	redis.nodeMode().deployMode().setDefaultValue()
}

func (redis *redis) setDefaultValue() *redis {

	if redis.mode == single {
		if !redis.remoteDeploy && redisListenAddress == "0.0.0.0" {
			redis.listenAddress = util.ExecuteCmdAcceptResult(constant.LocalIPCmd)
		}
	} else if redis.mode == cluster {
		if redis.localCluster && redisListenAddress == "0.0.0.0" {
			redis.listenAddress = util.ExecuteCmdAcceptResult(constant.LocalIPCmd)
		}
		if len(redis.clusterNodes) == 1 {
			redis.listenAddress = redis.clusterNodes[0].Host
		}
	}

	if redisListenAddress != "0.0.0.0" {
		redis.listenAddress = redisListenAddress
	}

	return redis
}
func (redis *redis) require() {
	// todo 检测参数合法性逻辑
	if redis.localCluster && redis.listenAddress == "0.0.0.0" {
		log.Fatal("当前模式必须指定--listen-address参数...")
	}
}

// 解析是否远程安装
func (redis *redis) deployMode() *redis {
	if len(redis.clusterNodes) > 0 {
		redis.remoteDeploy = true
	}

	return redis
}

// 解析节点数量
func (redis *redis) nodeMode() *redis {
	fmt.Println(len(redis.clusterNodes))

	if redisDeployMode == single {
		redis.mode = single
		if len(redis.clusterNodes) > 0 {
			redis.remoteDeploy = true
		} else {
			redis.remoteDeploy = false
		}
	} else if redisDeployMode == cluster {
		redis.mode = cluster
		if len(redis.clusterNodes) == 0 {
			redis.remoteDeploy = false
			redis.localCluster = true
			redis.clusterNodesNum = 1
		} else if len(redis.clusterNodes) != 0 {
			redis.clusterNodesNum = len(redis.clusterNodes)
			redis.remoteDeploy = true
		}
	} else {
		log.Fatal("请检查server list配置文件及参数...")
	}

	redis.banner(nil, fmt.Sprintf("redis部署模式：%s 节点数为： %d", redis.mode, len(redis.clusterNodes)), "")
	return redis
}

// 解析redis集群节点信息
func (redis *redis) cluterNodesInfo() {
	if _, err := os.Stat(serverListFile); err == nil {
		redis.banner(nil, fmt.Sprintf("解析%s文件", serverListFile), constant.Local)
		redis.clusterNodes = util.ParseServerList(serverListFile).RedisServerList
	} else {
		log.Fatal(err.Error())
	}
}

// 初始化集群
func (redis *redis) initializeCluster() {
	if redis.mode == cluster {
		switch redis.remoteDeploy {
		case false:
			redis.banner(nil, initRedisCluster, constant.LoopbackAddress)
			redis.shell(initRedisCluster, redis.initializeClusterCmd())
		case true:
			redis.banner(nil, initRedisCluster, redis.clusterNodes[0].Host)
			redis.clusterNodes[0].ExecuteOriginCmd(redis.initializeClusterCmd())
		}
	}
}

func (redis *redis) initializeClusterCmd() (cmd string) {
	cmd = "echo \"yes\" | redis-cli --cluster create"
	ports := []string{"26379", "26380", "26381", "26382", "26383", "26384"}

	var j int
	if redis.localCluster {
		j = 6
	} else {
		j = 6 / len(redis.clusterNodes)
	}

	if redis.localCluster {
		for i := 0; i < j; i++ {
			cmd += fmt.Sprintf(" %s:%s ", redisListenAddress, ports[i])
			// 赋值节点地址信息
			redis.connectAddress = append(redis.connectAddress, fmt.Sprintf("%s:%s", redis.listenAddress, ports[i]))
		}
	} else {
		for _, v := range redis.clusterNodes {
			for i := 0; i < j; i++ {
				cmd += fmt.Sprintf(" %s:%s ", v.Host, ports[i])
				// 赋值节点地址信息
				redis.connectAddress = append(redis.connectAddress, fmt.Sprintf("%s:%s", v.Host, ports[i]))
			}
		}
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

	if redis.mode == single {
		cmd = fmt.Sprintf("%s/redis-server %s/redis.conf",
			redisBinaryPath, redisConfigDir)
	} else if redis.mode == cluster {
		i := 6 / redis.clusterNodesNum
		for j := 0; j < i; j++ {
			cmd += fmt.Sprintf("redis-server %s/%s.conf;\n", redisConfigDir, ports[j])
		}
	}
	return cmd
}

func (redis *redis) banner(customBanner []string, msg string, address string) {
	var location string

	if redis.remoteDeploy {
		location = constant.Remote
	} else {
		location = constant.Local
		address = constant.LoopbackAddress
	}
	banners := []string{constant.Redis, location, address}
	for _, v := range customBanner {
		banners = append(banners, v)
	}
	fmt.Printf("%s %s...\n",
		util.PrintOrangeMulti(banners), msg)
}

// 配置systemd service
func (redis *redis) service() {
	redis.shell(configService, redis.serviceCmd())
}

func (redis *redis) serviceCmd() (cmd string) {
	if redis.mode == single {
		cmd = fmt.Sprintf("%s %s;sed -i \"/redis-server/d\" %s;echo \"%s/redis-server %s/redis.conf\" >> %s",
			constant.ChmodX, constant.EtcRcLocal, constant.EtcRcLocal, redisBinaryPath, redisConfigDir, constant.EtcRcLocal)
	}
	if redis.mode == cluster {
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

func (redis *redis) generateConfigFile() {

	var j int
	ports := []string{"26379", "26380", "26381", "26382", "26383", "26384"}

	if redis.mode == single && !redis.remoteDeploy {
		redis.banner(nil, generateRedisConfigFile, redisListenAddress)
		util.OverwriteContent(fmt.Sprintf("%s/redis.conf", redisConfigDir), redis.config("6379"))

	} else if redis.mode == single && redis.remoteDeploy {
		redis.banner(nil, generateRedisConfigFile, redisListenAddress)
		util.RemoteWriteFile(fmt.Sprintf("%s/redis.conf", redisConfigDir), []byte(redis.config("6379")), redis.clusterNodes[0])
	} else if redis.mode == cluster && redis.remoteDeploy {
		j = 6 / redis.clusterNodesNum
	} else if redis.mode == cluster && !redis.remoteDeploy {
		j = 6
	}

	for i := 0; i < j; i++ {
		if redis.remoteDeploy {
			for _, v := range redis.clusterNodes {
				configFilePath := fmt.Sprintf("%s/%s.conf", redisConfigDir, ports[i])
				redis.banner(nil, fmt.Sprintf("%s: %s", generateRedisConfigFile, configFilePath), v.Host)
				util.RemoteWriteFile(configFilePath, []byte(redis.config(ports[i])), v)
			}
		} else {
			redis.banner(nil, fmt.Sprintf("%s: %s/%s.conf...", generateRedisConfigFile, redisConfigDir, ports[i]), constant.LoopbackAddress)
			util.OverwriteContent(fmt.Sprintf("%s/%s.conf", redisConfigDir, ports[i]), redis.config(ports[i]))
		}

	}
}

// 编译
func (redis *redis) compile() {
	directoryName := util.CutCharacter(strings.TrimSuffix(sourceFilePath, ".tar.gz"), []string{"./", ""})
	compileCmd := fmt.Sprintf("tar zxvf %s >/dev/null && cd %s && sed -i \"s#\\$(PREFIX)/bin#%s#g\" src/Makefile && make -j $(nproc) && make install && cd ~",
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
	if redis.mode == single {
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
		redis.banner(nil, compileRedisEnvDetection, constant.Local)
		if !util.ExecuteIgnoreStd(search) {
			if !util.ExecuteIgnoreStd(install) {
				fmt.Printf("%s %s...\n",
					util.PrintRedMulti([]string{constant.Error, constant.Local, constant.LoopbackAddress}),
					dependenceDetectionNotPass)
			}
		}
	} else {
		for _, v := range redis.clusterNodes {
			redis.banner(nil, compileRedisEnvDetection, v.Host)
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

	if redis.mode == cluster {
		i = 6 / redis.clusterNodesNum
	} else if redis.mode == single {
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
	return fmt.Sprintf("pkill -9 redis;rm -f %s/nodes-263* %s/263*.conf %s/* %s/*",
		redisConfigDir, redisConfigDir, redisDataDir, redisLogDir)
}

func (redis *redis) status() {
	if redis.mode == cluster && redis.remoteDeploy {
		ssh := redis.clusterNodes[0]
		redis.banner(nil, clusterStatus, ssh.Host)
		ssh.RemoteShellPrint(fmt.Sprintf("%s/redis-cli -p 26379 -a %s -c cluster nodes",
			redisBinaryPath, redisPassword))
	} else if redis.mode == cluster && !redis.remoteDeploy {
		redis.banner(nil, clusterStatus, redis.listenAddress)
		util.ExecuteCmdAcceptResult(fmt.Sprintf("%s/redis-cli -p 26379 -a %s -c cluster nodes",
			redisBinaryPath, redisPassword))
	}
}

func (redis *redis) connection() {
	cmd := fmt.Sprintf("%s=%s\n", springRedisPassword, redisPassword)
	var address string
	if redis.mode == cluster {
		redis.banner(nil, redisConnectionInfo, constant.Local)
		for _, v := range redis.connectAddress {
			address += fmt.Sprintf("%s,", v)
		}
		cmd += fmt.Sprintf("%s=%s", springRedisClusterNodes, address)
		fmt.Println(strings.TrimSuffix(cmd, ","))
	} else if redis.mode == single {
		redis.banner(nil, redisConnectionInfo, redis.listenAddress)
		cmd += fmt.Sprintf("%s=%s\n%s=%d\n", springRedisHost, redis.listenAddress, springRedisPort, 6379)
		fmt.Println(cmd)
	}
}

func (redis *redis) shell(msg string, cmd string) {

	if !redis.remoteDeploy {
		redis.banner(nil, msg, constant.LoopbackAddress)
		util.ExecuteCmdAcceptResult(cmd)
	} else {
		wg := sync.WaitGroup{}
		wg.Add(len(redis.clusterNodes))
		for _, v := range redis.clusterNodes {
			redis.banner([]string{parallel}, msg, v.Host)
			go v.RemoteShellParallel(cmd, &wg)
		}
		wg.Wait()
	}
}
