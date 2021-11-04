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

const PruneDockerShell = `
systemctl disable docker --now || true
rm -rf /etc/docker || true
`

type DockerExternalConfig struct {
	Server []struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password int    `yaml:"password"`
		Port     int    `yaml:"port"`
	} `yaml:"server"`
	Excludes []string `yaml:"excludes"`
	Docker   struct {
		Package            string   `yaml:"package"`
		PreserveDir        string   `yaml:"preserveDir"`
		InsecureRegistries []string `yaml:"insecureRegistries"`
	} `yaml:"docker"`
}

// DockerInternalConfig 内部对象
type DockerInternalConfig struct {
	Servers            []runner.ServerInternal
	Package            string
	PreserveDir        string
	Logger             *logrus.Logger
	ConfigContent      []byte
	Executor           runner.ExecutorInternal
	IgnoreErr          bool // UnitTest
	InsecureRegistries []string
	BootCommand        string
	unitTest           bool
}

// Docker 安装部署docker-ce
func Docker(item command.OperationItem) (err command.RunErr) {
	//defer errors.IgnoreErrorFromCaller(3, "testing.tRunner", &err.Err)
	config := &DockerInternalConfig{
		Logger:        item.Logger,
		ConfigContent: item.B,
		unitTest:      item.UnitTest,
	}

	return install(config)
}

func (config *DockerInternalConfig) Parse() command.RunErr {

	config.Logger = log.SetDefault(config.Logger)

	obj := DockerExternalConfig{}
	config.Logger.Info("解析docker安装配置")
	if err := yaml.Unmarshal(config.ConfigContent, &obj); err != nil {
		return command.RunErr{Err: err}
	}

	// 深拷贝属性
	config.Package = obj.Docker.Package
	config.InsecureRegistries = obj.Docker.InsecureRegistries
	config.PreserveDir = obj.Docker.PreserveDir

	servers, err := runner.ParseServerList(config.ConfigContent, config.Logger)
	if err != nil {
		return command.RunErr{Err: fmt.Errorf("[docker] 反序列化主机列表失败 -> %s", err)}
	}

	config.Logger.Debugf("主机列表为：%v", servers)
	config.Servers = servers

	return command.RunErr{}
}

func (config *DockerInternalConfig) SetValue() command.RunErr {
	return command.RunErr{}
}

// Detect 调用运行时检测依赖
func (config *DockerInternalConfig) Detect() (err command.RunErr) {
	// todo:检测目录合法性
	return command.RunErr{}
}

// Prune 清理历史文件
func (config *DockerInternalConfig) Prune() (err command.RunErr) {

	defer func() {
		if err.Err != nil {
			err.Msg = "清理历史文件失败"
			if config.unitTest {
				err.Err = nil
			}
		}
	}()

	config.Logger = log.SetDefault(config.Logger)
	config.Logger.Infoln("清理docker历史文件...")
	exec := runner.ExecutorInternal{
		Servers:        config.Servers,
		Script:         PruneDockerShell,
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
func (config *DockerInternalConfig) HandPackage() (err command.RunErr) {

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

	config.Logger.Infoln("分发docker安装包完毕...")
	return command.RunErr{}
}

// Compile 编译
func (config *DockerInternalConfig) Compile() (err command.RunErr) {

	return command.RunErr{}
}

// SetUpRuntime docker运行时配置
func (config *DockerInternalConfig) SetUpRuntime() (err command.RunErr) {
	return command.RunErr{}
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
