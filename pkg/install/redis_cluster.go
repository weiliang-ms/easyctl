package install

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/install/tmpl"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/errors"
	"github.com/weiliang-ms/easyctl/pkg/util/log"
	strings2 "github.com/weiliang-ms/easyctl/pkg/util/strings"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

// RedisClusterConfig redis安装配置反序列化对象
type RedisClusterConfig struct {
	RedisCluster struct {
		Password    string           `yaml:"password"`
		ClusterType RedisClusterType `yaml:"cluster-type"`
		Package     string           `yaml:"package"`
		ListenPorts []int            `yaml:"listenPorts"`
	} `yaml:"redis-cluster"`
}

// 内部对象
type redisClusterConfig struct {
	Servers       []runner.ServerInternal
	Password      string
	CluterType    RedisClusterType
	ListenPorts   []int
	Package       string
	Logger        *logrus.Logger
	ConfigContent []byte
	ConfigItem    RedisClusterConfig
	Executor      runner.ExecutorInternal
	IgnoreErr     bool     // UnitTest
	PortsNeedOpen []int    // 需要放开的防火墙策略
	EndpointList  []string // 节点列表
	BootCommand   string
	unitTest      bool
}

// RedisClusterType redis cluster部署模式
type RedisClusterType int8

const (
	local RedisClusterType = iota
	threeNodesThreeShards
	sixNodesThreeShards
)

var defaultPorts = []int{26379, 26380, 26381, 26382, 26383, 26384}

const PruneRedisShell = `
pkill -9 redis || true
userdel -r redis || true
rm -rf /etc/redis || true
rm -rf /tmp/redis* || true
rm -rf /usr/local/bin/redis-* || true
rm -rf /var/lib/redis || true
rm -rf /var/log/redis || true
rm -f /var/run/redis* || true
rm -f /etc/init.d/redis* || true
`

// todo: io实现
const setUpRedisRuntimeShell = `
useradd redis -s /sbin/nologin -M || true
mkdir -p /var/lib/redis
mkdir -p /var/log/redis
mkdir -p /etc/redis
chown redis:redis /var/lib/redis
chown redis:redis /var/log/redis
chown redis:redis /etc/redis
`

// RedisCluster 部署redis集群
func RedisCluster(item command.OperationItem) command.RunErr {
	config := &redisClusterConfig{
		Logger:        item.Logger,
		ConfigContent: item.B,
	}

	config.unitTest = item.UnitTest

	return install(config)
}

func (config *redisClusterConfig) Parse() command.RunErr {

	if config.Logger == nil {
		config.Logger = logrus.New()
	}

	config.ConfigItem = RedisClusterConfig{}
	config.Logger.Info("解析redis cluster安装配置")
	if err := yaml.Unmarshal(config.ConfigContent, &config.ConfigItem); err != nil {
		return command.RunErr{Err: err}
	}

	// 深拷贝属性
	config.Package = config.ConfigItem.RedisCluster.Package
	config.CluterType = config.ConfigItem.RedisCluster.ClusterType
	config.Password = config.ConfigItem.RedisCluster.Password
	config.ListenPorts = config.ConfigItem.RedisCluster.ListenPorts

	servers, err := runner.ParseServerList(config.ConfigContent, config.Logger)
	if err != nil {
		return command.RunErr{Err: fmt.Errorf("[redis-cluster] 反序列化主机列表失败 -> %s", err)}
	}

	config.Logger.Debugf("主机列表为：%v", servers)
	config.Logger.Debugf("开始安装redis集群，集群模式为: %d", config.CluterType)
	config.Servers = servers

	return command.RunErr{}
}

func (config *redisClusterConfig) SetValue() command.RunErr {
	var ports []int
	var endpointList []string

	setLocalPorts := config.CluterType == local && len(config.ListenPorts) != 6
	setThreeNodesPorts := config.CluterType == threeNodesThreeShards && len(config.ListenPorts) != 2
	setSixNodesPorts := config.CluterType == sixNodesThreeShards && len(config.ListenPorts) != 1

	if setLocalPorts {
		ports = defaultPorts[:6]
	}

	if setThreeNodesPorts {
		ports = defaultPorts[:2]
	}

	if setSixNodesPorts {
		ports = defaultPorts[:1]
	}

	for _, i := range config.ListenPorts {
		ports = append(ports, i)
		ports = append(ports, i+10000)
		for _, v := range config.Servers {
			endpointList = append(endpointList, fmt.Sprintf("%s:%d", v.Host, i))
		}
	}

	config.PortsNeedOpen = ports
	config.EndpointList = endpointList

	return command.RunErr{}
}

// Detect 调用运行时检测依赖
func (config *redisClusterConfig) Detect() (err command.RunErr) {
	defer func() {
		if err.Err != nil && config.unitTest {
			err.Err = nil
		}
	}()

	const check = "gcc -v"

	if _, err := os.Stat(config.Package); err != nil {
		return command.RunErr{Err: errors.FileNotFoundErr(config.Package)}
	}

	if config.CluterType == threeNodesThreeShards && len(config.Servers) != 3 {
		return command.RunErr{Err: errors.NumNotEqualErr("节点", 3, len(config.Servers))}
	}

	if config.CluterType == sixNodesThreeShards && len(config.Servers) != 6 {
		return command.RunErr{Err: errors.NumNotEqualErr("节点", 6, len(config.Servers))}
	}

	config.Logger.Infoln("检测依赖环境...")

	config.Executor = runner.ExecutorInternal{
		Servers: config.Servers,
		Script:  check,
		Logger:  config.Logger,
	}

	if config.CluterType == local {
		config.Servers = config.Servers[:1]
	}

	for v := range config.Executor.ParallelRun() {
		if config.IgnoreErr {
			break
		} else if v.Err != nil {
			return command.RunErr{Err: fmt.Errorf("%s 依赖检测失败 -> %s", v.Host, v.Err)}
		}
	}

	return command.RunErr{}
}

// Prune 清理历史文件
func (config *redisClusterConfig) Prune() (err command.RunErr) {
	defer func() {
		if err.Err != nil && config.unitTest {
			err.Err = nil
		}
	}()

	config.Logger.Infoln("清理redis历史文件...")
	if config.CluterType == local {
		config.Servers = config.Servers[:1]
	}

	exec := runner.ExecutorInternal{
		Servers:        config.Servers,
		Script:         PruneRedisShell,
		Logger:         config.Logger,
		OutPutRealTime: true,
	}

	ch := exec.ParallelRun()
	for v := range ch {
		if config.IgnoreErr {
			break
		} else if v.Err != nil {
			return command.RunErr{Err: fmt.Errorf("[%s] 执行清理指令失败 %s", v.Host, v.Err)}
		}
	}

	return command.RunErr{}
}

// HandPackage 分发安装包
func (config *redisClusterConfig) HandPackage() (err command.RunErr) {

	defer func() {
		if err.Err != nil && config.unitTest {
			err.Err = nil
		}
	}()

	if config.CluterType == local {
		return command.RunErr{Err: os.Rename(config.Package, fmt.Sprintf("/tmp/%s",
			strings2.SubFileName(config.Package)))}
	}

	config.Logger.Infoln("分发package...")
	ch := runner.ParallelScp(runner.ScpItem{
		Servers: config.Servers,
		SrcPath: config.Package,
		DstPath: fmt.Sprintf("/tmp/%s",
			strings2.SubFileName(config.Package)),
		Mode:   0755,
		Logger: config.Logger,
	})

	for v := range ch {
		if config.IgnoreErr {
			break
		} else if v != nil {
			return command.RunErr{Err: v}
		}
	}

	config.Logger.Infoln("分发redis安装包完毕...")
	return command.RunErr{}
}

// Install 编译安装
func (config *redisClusterConfig) Install() (err command.RunErr) {

	defer func() {
		if err.Err != nil && config.unitTest {
			err.Err = nil
		}
	}()

	log.SetDefault(config.Logger)

	config.Logger.Infoln("开始编译redis")
	compileCmd, _ := tmplutil.Render(tmpl.RedisCompileTmpl, tmplutil.TmplRenderData{
		"PackageName": strings2.SubFileName(config.Package),
	})

	return command.RunErr{Err: config.run(compileCmd)}
}

// SetUpRuntime redis运行时配置
func (config *redisClusterConfig) SetUpRuntime() (err command.RunErr) {
	defer func() {
		if err.Err != nil && config.unitTest {
			err.Err = nil
		}
	}()
	config.Logger.Info("配置redis运行时环境")
	return command.RunErr{Err: config.run(setUpRedisRuntimeShell)}
}

func (config *redisClusterConfig) Config() (err command.RunErr) {

	defer func() {
		if err.Err != nil && config.unitTest {
			err.Err = nil
		}
	}()

	config.Logger = log.SetDefault(config.Logger)

	config.Logger.Info("生成配置文件")
	// local

	// todo: 考虑io替代shell
	generateConfigShell, _ := tmplutil.Render(tmpl.RedisConfigTmpl, tmplutil.TmplRenderData{
		"Ports":          config.ListenPorts,
		"Password":       config.Password,
		"ClusterEnabled": true,
	})

	return command.RunErr{Err: config.run(generateConfigShell)}
}

func (config *redisClusterConfig) SetService() (err command.RunErr) {
	defer func() {
		if err.Err != nil && config.unitTest {
			err.Err = nil
		}
	}()

	config.Logger.Info("配置开机自启动redis")

	config.Logger.Debugf("ports: %v", config.ListenPorts)

	// todo: 考虑io替代shell
	setServiceShell, _ := tmplutil.Render(tmpl.SetRedisServiceTmpl, tmplutil.TmplRenderData{
		"Ports":    config.ListenPorts,
		"Password": config.Password,
	})

	return command.RunErr{Err: config.run(setServiceShell)}
}

func (config *redisClusterConfig) Boot() (err command.RunErr) {

	defer func() {
		if err.Err != nil && config.unitTest {
			err.Err = nil
		}
	}()

	config.Logger.Info("启动redis")

	config.Logger.Debugf("ports: %v", config.ListenPorts)

	// todo: 考虑io替代shell
	bootRedisShell, _ := tmplutil.Render(tmpl.RedisBootTmpl, tmplutil.TmplRenderData{
		"Ports":    config.ListenPorts,
		"Password": config.Password,
	})

	config.BootCommand = bootRedisShell

	return command.RunErr{Err: config.run(bootRedisShell)}
}

func (config *redisClusterConfig) CloseFirewall() (err command.RunErr) {
	defer func() {
		if err.Err != nil && config.unitTest {
			err.Err = nil
		}
	}()

	config.Logger.Info("开放防火墙端口")

	script, _ := tmplutil.Render(tmpl.OpenFirewallPortTmpl, tmplutil.TmplRenderData{
		"Ports": config.PortsNeedOpen,
	})

	return command.RunErr{Err: config.run(script)}
}

func (config *redisClusterConfig) Init() (err command.RunErr) {
	defer func() {
		if err.Err != nil && config.unitTest {
			err.Err = nil
		}
	}()

	config.Logger.Info("初始化redis集群")

	script, _ := tmplutil.Render(tmpl.InitClusterTmpl, tmplutil.TmplRenderData{
		"EndpointList": config.EndpointList,
		"Password":     config.Password,
	})

	if len(config.Servers) > 0 {
		config.Servers = config.Servers[:1]
	}

	return command.RunErr{Err: config.run(script)}
}

func (config *redisClusterConfig) Print() command.RunErr {
	config.Logger.Info("redis集群安装完毕,相关信息如下：")
	var endpoint string
	for _, v := range config.EndpointList {
		endpoint = fmt.Sprintf("%s,%s", endpoint, v)
	}
	fmt.Printf("1.节点列表: %s\n"+
		"2.密码: %s\n"+
		"3.日志目录: /var/log/redis\n"+
		"4.数据目录: /var/data/redis\n"+
		"5.启动命令/节点: %s\n"+
		"6.二进制目录：/usr/local/bin/redis-*", strings.TrimPrefix(endpoint, ","), config.Password, config.BootCommand)
	return command.RunErr{}
}

func (config *redisClusterConfig) run(script string) error {

	if config.CluterType == local {
		if len(config.Servers) > 0 {
			config.Servers = config.Servers[:1]
		} else {
			return fmt.Errorf("单机集群：server不能为空也需指定")
		}
	}

	exec := runner.ExecutorInternal{
		Servers:        config.Servers,
		Script:         script,
		Logger:         config.Logger,
		OutPutRealTime: true,
	}

	ch := exec.ParallelRun()
	for v := range ch {
		if v.Err != nil {
			return v.Err
		}
	}

	return nil
}

func (config *redisClusterConfig) numPerNode() int {
	switch config.CluterType {
	case local:
		return 6
	case threeNodesThreeShards:
		return 2
	case sixNodesThreeShards:
		return 1
	default:
		return 0
	}
}
