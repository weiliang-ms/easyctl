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
systemctl stop docker.socket
systemctl disable docker --now || true
rm -rf /etc/docker || true
rm -f /usr/bin/containerd || true
rm -f /usr/bin/containerd-shim || true
rm -f /usr/bin/ctr || true
rm -f /usr/bin/docker || true
rm -f /usr/bin/docker-init || true
rm -f /usr/bin/docker-proxy || true
rm -f /usr/bin/dockerd || true
rm -f /usr/bin/runc || true
userdel -r docker || true
`

const bootDockerShell = `
setenforce 0
groupadd docker
useradd docker -g docker
systemctl daemon-reload
systemctl restart docker
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
		RegistryMirrors    []string `yaml:"registryMirrors"`
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
	Mirrors            []string
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
	config.Mirrors = obj.Docker.RegistryMirrors

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
	// todo:检测内核版本
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

	return command.RunErr{Err: config.run(PruneDockerShell)}
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
	if len(config.Servers) < 1 {
		re := runner.LocalRun(fmt.Sprintf("cp %s /tmp/%s", config.Package, config.Package), config.Logger)
		if re.Err != nil {
			fmt.Printf("%v", re)
			return command.RunErr{Err: err}
		}
	} else {
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
	}

	config.Logger.Infoln("分发docker安装包完毕...")
	return command.RunErr{}
}

// Install 编译&安装
func (config *DockerInternalConfig) Install() (err command.RunErr) {
	shell, _ := tmplutil.Render(tmpl.InstallDockerTmpl, tmplutil.TmplRenderData{
		"Package": config.Package,
	})

	return command.RunErr{Err: config.run(shell)}
}

// SetUpRuntime docker运行时配置
func (config *DockerInternalConfig) SetUpRuntime() (err command.RunErr) {
	return command.RunErr{Err: config.run(fmt.Sprintf("mkdir -p %s", config.PreserveDir))}
}

// Config 配置/etc/docker/daemon.json
func (config *DockerInternalConfig) Config() (err command.RunErr) {
	defer func() {
		if err.Err != nil {
			err.Msg = "配置docker异常"
			if config.unitTest {
				err.Err = nil
			}
		}
	}()

	config.Logger = log.SetDefault(config.Logger)

	config.Logger.Info("生成配置文件")

	var Mirrors, InsecureRegistries string
	if config.Mirrors != nil {
		var mirrors []string
		for _, mirror := range config.Mirrors {
			mirrors = append(mirrors, fmt.Sprintf("\"%s\"", mirror))
		}
		Mirrors = strings.Join(mirrors, ", ")
	}
	if config.InsecureRegistries != nil {
		var registries []string
		for _, repostry := range config.InsecureRegistries {
			registries = append(registries, fmt.Sprintf("\"%s\"", repostry))
		}
		InsecureRegistries = strings.Join(registries, ", ")
	}
	generateConfigShell, _ := tmplutil.Render(tmpl.DockerConfigTmpl, tmplutil.TmplRenderData{
		"Mirrors":            Mirrors,
		"InsecureRegistries": InsecureRegistries,
		"DataPath":           config.PreserveDir,
	})

	return command.RunErr{Err: config.run(generateConfigShell)}
}

func (config *DockerInternalConfig) SetService() (err command.RunErr) {
	defer func() {
		if err.Err != nil {
			err.Msg = "配置docker开机启动异常"
			if config.unitTest {
				err.Err = nil
			}
		}
	}()

	config.Logger.Info("配置开机自启动docker")

	// todo: 考虑io替代shell
	setServiceShell, _ := tmplutil.Render(tmpl.SetDockerServiceTmpl, tmplutil.TmplRenderData{})

	return command.RunErr{Err: config.run(setServiceShell)}
}

func (config *DockerInternalConfig) Boot() (err command.RunErr) {

	defer func() {
		if err.Err != nil {
			err.Msg = "启动docker异常"
			if config.unitTest {
				err.Err = nil
			}
		}
	}()

	config.Logger.Info("启动docker")

	return command.RunErr{Err: config.run(bootDockerShell)}
}

func (config *DockerInternalConfig) CloseFirewall() (err command.RunErr) {
	return command.RunErr{}
}

func (config *DockerInternalConfig) Init() (err command.RunErr) {
	return command.RunErr{}
}

func (config *DockerInternalConfig) Print() command.RunErr {
	config.Logger.Info("docker安装完毕")
	return command.RunErr{}
}

func (config *DockerInternalConfig) run(script string) error {

	if len(config.Servers) < 1 {
		runner.LocalRun(script, config.Logger)
	} else {
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
	}

	return nil
}
