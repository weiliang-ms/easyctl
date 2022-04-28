package docker

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/install"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/log"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
	"gopkg.in/yaml.v2"
	"strings"
	"sync"
	"time"
)

// PruneDockerShell todo: optimize
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

type ManifestConfig struct {
	Docker struct {
		Package            string   `yaml:"package"`
		PreserveDir        string   `yaml:"preserveDir"`
		InsecureRegistries []string `yaml:"insecureRegistries"`
		RegistryMirrors    []string `yaml:"registryMirrors"`
	} `yaml:"docker"`
}

// Manager 内部对象
type Manager struct {
	Servers            []runner.ServerInternal
	Package            string
	PreserveDir        string
	Logger             *logrus.Logger
	ConfigContent      []byte
	Timeout            time.Duration
	Executor           runner.ExecutorInternal
	IgnoreErr          bool // UnitTest
	InsecureRegistries []string
	Mirrors            []string
	BootCommand        string
	unitTest           bool
	Local              bool // 本地安装
	Handler            HandlerInterface
}

// Install 安装部署docker-ce
func Install(item command.OperationItem) command.RunErr {
	//defer errors.IgnoreErrorFromCaller(3, "testing.tRunner", &err.Err)
	m := &Manager{
		Logger:        item.Logger,
		ConfigContent: item.B,
		unitTest:      item.UnitTest,
		Local:         item.LocalRun,
		Timeout:       item.SSHTimeout,
	}

	m.Handler = getHandlerInterface(item.Interface)

	result, err := m.Parse()
	if err != nil {
		return command.RunErr{Err: err}
	}

	f := func(s runner.ServerInternal) command.RunErr {

		if err := m.Detect(s); err.Err != nil {
			return err
		}

		if err := m.Prune(s); err.Err != nil {
			return command.RunErr{Err: err}
		}

		if err := m.HandPackage(s); err.Err != nil {
			return err
		}

		if err := m.Install(s); err.Err != nil {
			return command.RunErr{Err: err}
		}

		if err := m.SetUpRuntime(s); err.Err != nil {
			return command.RunErr{Err: err}
		}

		if err := m.SetConfig(s); err.Err != nil {
			return command.RunErr{Err: err}
		}

		if err := m.SetSystemd(s); err.Err != nil {
			return command.RunErr{Err: err}
		}

		if err := m.Boot(s); err.Err != nil {
			return command.RunErr{Err: err}
		}

		return command.RunErr{}
	}

	if !m.Local {
		wg := sync.WaitGroup{}
		wg.Add(len(result.Servers))
		ch := make(chan command.RunErr, len(m.Servers))

		for _, v := range result.Servers {
			go func(server runner.ServerInternal) {
				ch <- f(server)
				wg.Done()
			}(v)

		}

		wg.Wait()
		close(ch)

		for v := range ch {
			if v.Err != nil {
				return v
			}
		}

		return command.RunErr{}
	}

	return f(runner.ServerInternal{})
}

// Parse todo: return Config object
func (m *Manager) Parse() (*Manager, error) {

	// todo 反射通用类型初始化（New）
	m.Logger = log.SetDefault(m.Logger)

	obj := ManifestConfig{}
	m.Logger.Info("解析docker安装配置")
	if err := yaml.Unmarshal(m.ConfigContent, &obj); err != nil {
		return m, err
	}

	// 深拷贝属性
	m.Package = obj.Docker.Package
	m.InsecureRegistries = obj.Docker.InsecureRegistries
	m.PreserveDir = obj.Docker.PreserveDir
	m.Mirrors = obj.Docker.RegistryMirrors

	servers, err := runner.ParseServerList(m.ConfigContent, m.Logger)
	if err != nil {
		return m, install.ParseServerListErr{Err: err}
	}

	m.Servers = servers
	m.Logger.Debugf("主机列表为：%v", servers)
	return m, nil
}

/*
	todo: goroutine. don't panic err when set mock
*/

// Detect 调用运行时检测依赖
func (m *Manager) Detect(server runner.ServerInternal) (err command.RunErr) {
	// todo:检测目录合法性
	// todo:检测内核版本
	return command.RunErr{Err: m.Handler.Detect("", server, m.Local, m.Logger, m.Timeout)}
}

// Prune 清理历史文件
// todo: mock
func (m *Manager) Prune(server runner.ServerInternal) install.PruneErr {
	m.Logger = log.SetDefault(m.Logger)
	m.Logger.Info("清理docker历史文件...")
	return m.Handler.Prune(server, m.Local, m.Logger, m.Timeout)
}

//func (m *Manager) Prune(server runner.ServerInternal) (err command.RunErr) {
//	m.Logger = log.SetDefault(m.Logger)
//	m.Logger.Info("清理docker历史文件...")
//	return command.RunErr{Err: m.Handler.Prune(server, m.Local, m.Logger, m.Timeout)}
//}

// HandPackage 分发安装包
// todo: mock
func (m *Manager) HandPackage(server runner.ServerInternal) (err command.RunErr) {
	logger := log.SetDefault(m.Logger)
	logger.Info("分发package...")
	return command.RunErr{Err: m.Handler.HandPackage(server, m.Package, m.Local, logger, m.Timeout)}
}

// Install 安装
func (m *Manager) Install(server runner.ServerInternal) install.InstallErr {
	return m.Handler.Install(tmplutil.RenderPanicErr(InstallDockerTmpl, tmplutil.TmplRenderData{
		"Package": m.Package,
	}), server, m.Local, m.Logger, m.Timeout)
}

// SetUpRuntime docker运行时配置
func (m *Manager) SetUpRuntime(server runner.ServerInternal) install.SetUpRuntimeErr {
	return m.Handler.SetUpRuntime(
		fmt.Sprintf("mkdir -p %s", m.PreserveDir),
		server,
		m.Local,
		m.Logger,
		m.Timeout,
	)
}

// SetConfig 配置/etc/docker/daemon.json
func (m *Manager) SetConfig(server runner.ServerInternal) install.SetConfigErr {

	m.Logger = log.SetDefault(m.Logger)
	m.Logger.Info("生成配置文件")

	var Mirrors, InsecureRegistries string
	if m.Mirrors != nil {
		var mirrors []string
		for _, mirror := range m.Mirrors {
			mirrors = append(mirrors, fmt.Sprintf("\"%s\"", mirror))
		}
		Mirrors = strings.Join(mirrors, ", ")
	}
	if m.InsecureRegistries != nil {
		var registries []string
		for _, repostry := range m.InsecureRegistries {
			registries = append(registries, fmt.Sprintf("\"%s\"", repostry))
		}
		InsecureRegistries = strings.Join(registries, ", ")
	}

	return m.Handler.SetConfig(
		tmplutil.RenderPanicErr(DockerConfigTmpl, tmplutil.TmplRenderData{
			"Mirrors":            Mirrors,
			"InsecureRegistries": InsecureRegistries,
			"DataPath":           m.PreserveDir,
		}),
		server, m.Local, m.Logger, m.Timeout,
	)
}

// SetSystemd 配置systemd服务
func (m *Manager) SetSystemd(server runner.ServerInternal) install.SetSystemdErr {
	m.Logger.Info("配置docker.service")
	return m.Handler.SetSystemd(
		tmplutil.RenderPanicErr(SetDockerServiceTmpl, tmplutil.TmplRenderData{}),
		server, m.Local, m.Logger, m.Timeout)
}

// Boot 启动docker服务
func (m *Manager) Boot(server runner.ServerInternal) install.BootErr {
	m.Logger.Info("启动docker")
	return m.Handler.Boot(server, m.Local, m.Logger, m.Timeout)
}

func getHandlerInterface(i interface{}) HandlerInterface {
	handlerInterface, _ := i.(HandlerInterface)
	if handlerInterface == nil {
		return new(Handler)
	}
	return handlerInterface
}
