package install

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/install/tmpl"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/log"
	strings2 "github.com/weiliang-ms/easyctl/pkg/util/strings"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
	"gopkg.in/yaml.v2"
	"strings"
)

// RedisExternalConfig redis安装配置反序列化对象
type RedisExternalConfig struct {
	Redis struct {
		Password string `yaml:"password"`
		Port     int    `yaml:"port"`
		Package  string `yaml:"package"`
	} `yaml:"redis"`
}

// RedisInternalConfig 内部对象
type RedisInternalConfig struct {
	Servers       []runner.ServerInternal
	Password      string
	Ports         []int
	Package       string
	Logger        *logrus.Logger
	ConfigContent []byte
	Executor      runner.ExecutorInternal
	IgnoreErr     bool     // UnitTest
	PortsNeedOpen []int    // 需要放开的防火墙策略
	EndpointList  []string // 节点列表
	BootCommand   string
	unitTest      bool
}

// Redis 部署redis集群
func Redis(item command.OperationItem) (err command.RunErr) {
	//defer errors.IgnoreErrorFromCaller(3, "testing.tRunner", &err.Err)
	config := &RedisInternalConfig{
		Logger:        item.Logger,
		ConfigContent: item.B,
		unitTest:      item.UnitTest,
	}

	return install(config)
}

func (config *RedisInternalConfig) Parse() command.RunErr {

	config.Logger = log.SetDefault(config.Logger)

	obj := RedisExternalConfig{}
	config.Logger.Info("解析redis安装配置")
	if err := yaml.Unmarshal(config.ConfigContent, &obj); err != nil {
		return command.RunErr{Err: err}
	}

	// 深拷贝属性
	config.Package = obj.Redis.Package
	config.Ports = []int{obj.Redis.Port}
	config.Password = obj.Redis.Password

	servers, err := runner.ParseServerList(config.ConfigContent, config.Logger)
	if err != nil {
		return command.RunErr{Err: fmt.Errorf("[redis] 反序列化主机列表失败 -> %s", err)}
	}

	config.Logger.Debugf("主机列表为：%v", servers)
	config.Servers = servers

	return command.RunErr{}
}

func (config *RedisInternalConfig) SetValue() command.RunErr {
	var ports []int
	var endpointList []string
	config.Logger = log.SetDefault(config.Logger)

	for _, p := range config.Ports {
		if p != 0 {
			ports = append(ports, p)
		}
	}

	if len(ports) == 0 {
		ports = append(ports, 6379)
	}

	config.Ports = ports

	for _, v := range config.Servers {
		for _, p := range config.Ports {
			endpointList = append(endpointList, fmt.Sprintf("%s:%d", v.Host, p))
		}
	}

	config.PortsNeedOpen = ports
	config.EndpointList = endpointList

	return command.RunErr{}
}

// Detect 调用运行时检测依赖
func (config *RedisInternalConfig) Detect() (err command.RunErr) {

	defer func() {
		if err.Err != nil && config.unitTest {
			err.Err = nil
		}
	}()

	// port合法性
	for _, v := range config.Ports {
		if v < 1 || v > 65535 {
			return command.RunErr{Err: fmt.Errorf("port: %d取值范围非法", v)}
		}
	}

	// 检测gcc
	const check = "gcc -v"
	config.Logger = log.SetDefault(config.Logger)
	config.Logger.Infoln("检测依赖环境...")
	config.Executor = runner.ExecutorInternal{
		Servers: config.Servers,
		Script:  check,
		Logger:  config.Logger,
	}

	for v := range config.Executor.ParallelRun() {
		if config.IgnoreErr {
			break
		} else if v.Err != nil {
			return command.RunErr{Err: fmt.Errorf("%s 依赖检测失败 -> %s", v.Host, v.Err), Msg: "依赖检测失败"}
		}
	}

	return command.RunErr{}
}

// Prune 清理历史文件
func (config *RedisInternalConfig) Prune() (err command.RunErr) {

	defer func() {
		if err.Err != nil {
			err.Msg = "清理历史文件失败"
			if config.unitTest {
				err.Err = nil
			}
		}
	}()

	config.Logger = log.SetDefault(config.Logger)
	config.Logger.Infoln("清理redis历史文件...")
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
			return command.RunErr{Err: fmt.Errorf("[%s] 执行清理指令失败 %s", v.Host, v.Err), Msg: "执行清理指令失败"}
		}
	}

	return command.RunErr{}
}

// HandPackage 分发安装包
func (config *RedisInternalConfig) HandPackage() (err command.RunErr) {

	defer func() {
		if err.Err != nil {
			err.Msg = "分发安装包失败"
			if config.unitTest {
				err.Err = nil
			}
		}
	}()

	config.Logger = log.SetDefault(config.Logger)
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

// Compile 编译
func (config *RedisInternalConfig) Compile() (err command.RunErr) {
	defer func() {
		if err.Err != nil && config.unitTest {
			err.Err = nil
			err.Msg = "编译redis异常"
		}
	}()

	config.Logger = log.SetDefault(config.Logger)

	config.Logger.Infoln("开始编译redis")
	compileCmd, _ := tmplutil.Render(tmpl.RedisCompileTmpl, tmplutil.TmplRenderData{
		"PackageName": strings2.SubFileName(config.Package),
	})

	return command.RunErr{Err: config.run(compileCmd)}
}

// SetUpRuntime redis运行时配置
func (config *RedisInternalConfig) SetUpRuntime() (err command.RunErr) {

	defer func() {
		if err.Err != nil {
			err.Msg = "配置运行时参数异常"
			if config.unitTest {
				err.Err = nil
			}
		}
	}()

	config.Logger.Info("配置redis运行时环境")
	return command.RunErr{Err: config.run(setUpRedisRuntimeShell)}
}

func (config *RedisInternalConfig) Config() (err command.RunErr) {
	defer func() {
		if err.Err != nil {
			err.Msg = "配置redis异常"
			if config.unitTest {
				err.Err = nil
			}
		}
	}()

	config.Logger = log.SetDefault(config.Logger)

	config.Logger.Info("生成配置文件")

	// todo: 考虑io替代shell
	generateConfigShell, _ := tmplutil.Render(tmpl.RedisConfigTmpl, tmplutil.TmplRenderData{
		"Ports":    config.Ports,
		"Password": config.Password,
	})

	return command.RunErr{Err: config.run(generateConfigShell)}
}

func (config *RedisInternalConfig) SetService() (err command.RunErr) {
	defer func() {
		if err.Err != nil {
			err.Msg = "配置redis开机启动异常"
			if config.unitTest {
				err.Err = nil
			}
		}
	}()

	config.Logger.Info("配置开机自启动redis")

	// local
	config.Logger.Debugf("ports: %v", config.Ports)

	// todo: 考虑io替代shell
	setServiceShell, _ := tmplutil.Render(tmpl.SetRedisServiceTmpl, tmplutil.TmplRenderData{
		"Ports":    config.Ports,
		"Password": config.Password,
	})

	return command.RunErr{Err: config.run(setServiceShell)}
}

func (config *RedisInternalConfig) Boot() (err command.RunErr) {

	defer func() {
		if err.Err != nil {
			err.Msg = "启动redis异常"
			if config.unitTest {
				err.Err = nil
			}
		}
	}()

	config.Logger.Info("启动redis")

	// local
	ports := config.Ports
	config.Logger.Debugf("ports: %v", ports)

	// todo: 考虑io替代shell
	bootRedisShell, _ := tmplutil.Render(tmpl.RedisBootTmpl, tmplutil.TmplRenderData{
		"Ports":    ports,
		"Password": config.Password,
	})

	config.BootCommand = bootRedisShell

	return command.RunErr{Err: config.run(bootRedisShell)}
}

func (config *RedisInternalConfig) CloseFirewall() (err command.RunErr) {
	defer func() {
		if err.Err != nil {
			err.Msg = "开放防火墙端口失败"
			if config.unitTest {
				err.Err = nil
			}
		}
	}()

	config.Logger.Info("开放防火墙端口")

	script, _ := tmplutil.Render(tmpl.OpenFirewallPortTmpl, tmplutil.TmplRenderData{
		"Ports": config.PortsNeedOpen,
	})

	return command.RunErr{Err: config.run(script)}
}

func (config *RedisInternalConfig) Init() (err command.RunErr) {
	return command.RunErr{}
}

func (config *RedisInternalConfig) Print() command.RunErr {
	config.Logger.Info("redis安装完毕,相关信息如下：")
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

func (config *RedisInternalConfig) run(script string) error {

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
