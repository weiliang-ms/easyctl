package install

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/install/tmpl"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/errors"
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
	} `yaml:"redis-cluster"`
}

// 内部对象
type redisClusterConfig struct {
	Servers       []runner.ServerInternal
	Password      string
	CluterType    RedisClusterType
	Package       string
	Logger        *logrus.Logger
	ConfigContent []byte
	ConfigItem    RedisClusterConfig
	Executor      runner.ExecutorInternal
	IgnoreErr     bool     // UnitTest
	PortsNeedOpen []int    // 需要放开的防火墙策略
	EndpointList  []string // 节点列表
	BootCommand   string
}

// RedisClusterType redis cluster部署模式
type RedisClusterType int8

const (
	local RedisClusterType = iota
	threeNodesThreeShards
	sixNodesThreeShards
)

var defaultPorts = []int{26379, 26380, 26381, 26382, 26383, 26384}

// todo: io实现
const pruneRedisShell = `
pkill -9 redis || true
userdel -r redis || true
rm -rf /etc/redis || true
rm -rf /tmp/redis* || true
rm -rf /usr/local/bin/redis-* || true
rm -rf /var/lib/redis || true
rm -rf /var/log/redis || true
rm -f /var/run/redis* || true
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
func RedisCluster(item command.OperationItem) (err error) {
	defer errors.IgnoreErrorFromCaller(3, "testing.tRunner", &err)
	config := &redisClusterConfig{
		Logger:        item.Logger,
		ConfigContent: item.B,
	}

	return install(config)
}

func (config *redisClusterConfig) Parse() error {

	if config.Logger == nil {
		config.Logger = logrus.New()
	}

	config.ConfigItem = RedisClusterConfig{}
	config.Logger.Info("解析redis cluster安装配置")
	if err := yaml.Unmarshal(config.ConfigContent, &config.ConfigItem); err != nil {
		return err
	}

	// 深拷贝属性
	config.Package = config.ConfigItem.RedisCluster.Package
	config.CluterType = config.ConfigItem.RedisCluster.ClusterType
	config.Password = config.ConfigItem.RedisCluster.Password

	servers, err := runner.ParseServerList(config.ConfigContent, config.Logger)
	if err != nil {
		return fmt.Errorf("[redis-cluster] 反序列化主机列表失败 -> %s", err)
	}

	config.Logger.Debugf("主机列表为：%v", servers)
	config.Logger.Debugf("开始安装redis集群，集群模式为: %d", config.CluterType)
	config.Servers = servers

	return nil
}

func (config *redisClusterConfig) SetValue() error {
	var ports []int
	var endpointList []string

	for _, i := range defaultPorts[:config.numPerNode()] {
		ports = append(ports, i)
		ports = append(ports, i+10000)
		for _, v := range config.Servers {
			endpointList = append(endpointList, fmt.Sprintf("%s:%d", v.Host, i))
		}
	}

	config.PortsNeedOpen = ports
	config.EndpointList = endpointList

	return nil
}

// Detect 调用运行时检测依赖
func (config *redisClusterConfig) Detect() (err error) {
	const check = "gcc -v"

	if _, err := os.Stat(config.Package); err != nil {
		return errors.FileNotFoundErr(config.Package)
	}

	if config.CluterType == threeNodesThreeShards && len(config.Servers) != 3 {
		return errors.NumNotEqualErr("节点", 3, len(config.Servers))
	}

	if config.CluterType == sixNodesThreeShards && len(config.Servers) != 6 {
		return errors.NumNotEqualErr("节点", 6, len(config.Servers))
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
			return fmt.Errorf("%s 依赖检测失败 -> %s", v.Host, v.Err)
		}
	}

	return nil
}

// Prune 清理历史文件
func (config *redisClusterConfig) Prune() error {

	config.Logger.Infoln("清理redis历史文件...")
	if config.CluterType == local {
		config.Servers = config.Servers[:1]
	}

	exec := runner.ExecutorInternal{
		Servers:        config.Servers,
		Script:         pruneRedisShell,
		Logger:         config.Logger,
		OutPutRealTime: true,
	}

	ch := exec.ParallelRun()
	for v := range ch {
		if config.IgnoreErr {
			break
		} else if v.Err != nil {
			return fmt.Errorf("[%s] 执行清理指令失败 %s", v.Host, v.Err)
		}
	}

	return nil
}

// HandPackage 分发安装包
func (config *redisClusterConfig) HandPackage() (err error) {

	if config.CluterType == local {
		defer errors.IgnoreErrorFromCaller(3, "testing.tRunner", &err)
		return os.Rename(config.Package, fmt.Sprintf("/tmp/%s",
			strings2.SubFileName(config.Package)))
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
			return v
		}
	}

	config.Logger.Infoln("分发redis安装包完毕...")
	return nil
}

// Compile 编译
func (config *redisClusterConfig) Compile() error {
	// todo: if nit set value
	if config.Logger == nil {
		config.Logger = logrus.New()
	}
	config.Logger.Infoln("开始编译redis")
	compileCmd, _ := tmplutil.Render(tmpl.RedisCompileTmpl, tmplutil.TmplRenderData{
		"PackageName": strings2.SubFileName(config.Package),
	})

	return config.run(compileCmd)
}

// SetUpRuntime redis运行时配置
func (config *redisClusterConfig) SetUpRuntime() error {
	config.Logger.Info("配置redis运行时环境")
	return config.run(setUpRedisRuntimeShell)
}

func (config *redisClusterConfig) Config() error {

	// todo: 非空校验
	if config.Logger == nil {
		config.Logger = logrus.New()
	}

	config.Logger.Info("生成配置文件")
	// local
	var ports []int
	if config.CluterType == local {
		ports = defaultPorts[:6]
	}

	if config.CluterType == threeNodesThreeShards {
		ports = defaultPorts[:2]
	}

	if config.CluterType == sixNodesThreeShards {
		ports = defaultPorts[:1]
	}

	// todo: 考虑io替代shell
	generateConfigShell, _ := tmplutil.Render(tmpl.RedisClusterConfigTmpl, tmplutil.TmplRenderData{
		"Ports":    ports,
		"Password": config.Password,
	})

	return config.run(generateConfigShell)
}

func (config *redisClusterConfig) SetService() error {
	config.Logger.Info("配置开机自启动redis")

	// local
	ports := defaultPorts[:config.numPerNode()]
	config.Logger.Debugf("ports: %v", ports)

	// todo: 考虑io替代shell
	setServiceShell, _ := tmplutil.Render(tmpl.SetRedisServiceTmpl, tmplutil.TmplRenderData{
		"Ports":    ports,
		"Password": config.Password,
	})

	return config.run(setServiceShell)
}

func (config *redisClusterConfig) Boot() error {

	config.Logger.Info("启动redis")

	// local
	ports := defaultPorts[:config.numPerNode()]
	config.Logger.Debugf("ports: %v", ports)

	// todo: 考虑io替代shell
	bootRedisShell, _ := tmplutil.Render(tmpl.RedisBootTmpl, tmplutil.TmplRenderData{
		"Ports":    ports,
		"Password": config.Password,
	})

	config.BootCommand = bootRedisShell

	return config.run(bootRedisShell)
}

func (config *redisClusterConfig) CloseFirewall() error {
	config.Logger.Info("开放防火墙端口")

	script, _ := tmplutil.Render(tmpl.OpenFirewallPortTmpl, tmplutil.TmplRenderData{
		"Ports": config.PortsNeedOpen,
	})

	return config.run(script)
}

func (config *redisClusterConfig) Init() error {
	config.Logger.Info("初始化redis集群")

	script, _ := tmplutil.Render(tmpl.InitClusterTmpl, tmplutil.TmplRenderData{
		"EndpointList": config.EndpointList,
		"Password":     config.Password,
	})

	if len(config.Servers) > 0 {
		config.Servers = config.Servers[:1]
	}

	return config.run(script)
}

func (config *redisClusterConfig) Print() error {
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
	return nil
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
