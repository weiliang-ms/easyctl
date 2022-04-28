package docker

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/weiliang-ms/easyctl/pkg/install"
	"github.com/weiliang-ms/easyctl/pkg/install/docker/mocks"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
	"gopkg.in/yaml.v2"
	"testing"
	"time"
)

var (
	mockHost    = "1.1.1.1"
	mockErr     = fmt.Errorf("mock err")
	mockTimeout = time.Millisecond

	mockBootErr = install.BootErr{
		Host: mockHost,
		Err:  mockErr,
	}

	mockSetSystemdErr = install.SetSystemdErr{
		Host: mockHost,
		Err:  mockErr,
	}

	mockSetConfigErr = install.SetConfigErr{
		Host: mockHost,
		Err:  mockErr,
	}

	mockSetRuntimeErr = install.SetUpRuntimeErr{
		Host: mockHost,
		Err:  mockErr,
	}

	mockInstallErr = install.InstallErr{
		Host: mockHost,
		Err:  mockErr,
	}

	mockPruneErr = install.PruneErr{
		Host: mockHost,
		Err:  mockErr,
	}

	mockTransferErr = install.TransferPackageErr{
		Host: mockHost,
		Path: mockPackagePath,
		Err:  mockErr,
	}

	mockDetectErr = install.DetectErr{
		Host: mockHost,
		Err:  mockErr,
	}

	mockMirrors                 = "\"docker.a.io\", \"docker.b.io\""
	mockMirrorsSlice            = []string{"docker.a.io", "docker.b.io"}
	mockInsecureRegistries      = "\"gcr.azk8s.cn\", \"quay.azk8s.cn\""
	mockInsecureRegistriesSlice = []string{"gcr.azk8s.cn", "quay.azk8s.cn"}
	mockLogger                  = logrus.New()

	mockPackagePath = "docker.tgz"

	mockSetRuntimeShell = "mkdir -p /data/docker"
	mockSetSystemdShell = tmplutil.RenderPanicErr(SetDockerServiceTmpl, tmplutil.TmplRenderData{})
	mockSetConfigShell  = tmplutil.RenderPanicErr(DockerConfigTmpl, tmplutil.TmplRenderData{
		"Mirrors":            mockMirrors,
		"InsecureRegistries": mockInsecureRegistries,
		"DataPath":           mockPreserveDir,
	})
	mockInstallShell = tmplutil.RenderPanicErr(InstallDockerTmpl, tmplutil.TmplRenderData{
		"Package": mockPackagePath,
	})
	mockPreserveDir = "/data/docker"
	mockLocal       = false
	mockServer      = runner.ServerInternal{
		Host:           "1.1.1.1",
		Port:           "22",
		UserName:       "root",
		Password:       "123456",
		PrivateKeyPath: "",
	}

	mockManager = Manager{
		Logger:             mockLogger,
		Package:            mockPackagePath,
		PreserveDir:        mockPreserveDir,
		InsecureRegistries: mockInsecureRegistriesSlice,
		Mirrors:            mockMirrorsSlice,
		Local:              false,
		Timeout:            mockTimeout,
	}

	mockInstallContent = `
server:
  - host: 1.1.1.1
    username: root
    password: 123456
    port: 22
docker:
 package: docker.tgz   # 二进制安装包目录
 preserveDir: /data/docker  # docker数据持久化目录
 insecureRegistries: # 非https仓库列表
   - gcr.azk8s.cn
   - quay.azk8s.cn
 registryMirrors:               # 镜像源
   - docker.a.io
   - docker.b.io
`
)

/*
	mock
*/

func Test_Boot_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface

	mockInterface.On("Boot", mockServer, mockLocal, mockLogger, mockTimeout).Return(install.BootErr{})
	require.Equal(t, install.BootErr{}, mockManager.Boot(mockServer))
}

func Test_Boot_Err_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface

	mockInterface.On("Boot", mockServer, mockLocal, mockLogger, mockTimeout).Return(mockBootErr)
	require.NotNil(t, mockManager.Boot(mockServer).Err)
}

func Test_SetSystemd_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface

	mockInterface.On("SetSystemd",
		mockSetSystemdShell, mockServer, mockLocal, mockLogger, mockTimeout).
		Return(install.SetSystemdErr{})

	require.Equal(t, install.SetSystemdErr{}, mockManager.SetSystemd(mockServer))
}

func Test_SetSystemd_Err_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface

	mockInterface.On("SetSystemd",
		mockSetSystemdShell, mockServer, mockLocal, mockLogger, mockTimeout).
		Return(mockSetSystemdErr)

	require.NotNil(t, mockManager.SetSystemd(mockServer).Err)
}

func Test_Config_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface

	mockInterface.On("SetConfig",
		mockSetConfigShell, mockServer, mockLocal, mockLogger, mockTimeout).
		Return(install.SetConfigErr{})

	require.Equal(t, install.SetConfigErr{}, mockManager.SetConfig(mockServer))
}

func Test_Config_Err_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface

	mockInterface.On("SetConfig",
		mockSetConfigShell, mockServer, mockLocal, mockLogger, mockTimeout).
		Return(mockSetConfigErr)

	require.NotNil(t, mockManager.SetConfig(mockServer).Err)
}

func Test_SetUpRuntime_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface

	mockInterface.On("SetUpRuntime",
		mockSetRuntimeShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.SetUpRuntimeErr{})

	require.Equal(t, install.SetUpRuntimeErr{}, mockManager.SetUpRuntime(mockServer))
}

func Test_SetUpRuntimeErr_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface

	mockInterface.On("SetUpRuntime",
		mockSetRuntimeShell, mockServer, mockLocal, mockLogger, mockTimeout).
		Return(mockSetRuntimeErr)

	require.NotNil(t, mockManager.SetUpRuntime(mockServer).Err)
}

func Test_InstallFunc_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface

	mockInterface.On("Install",
		mockInstallShell, mockServer, mockLocal, mockLogger, mockTimeout).
		Return(install.InstallErr{})

	require.Equal(t, install.InstallErr{}, mockManager.Install(mockServer))
}

func Test_InstallErr_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface

	mockInterface.On("Install",
		mockInstallShell, mockServer, mockLocal, mockLogger, mockTimeout).
		Return(mockInstallErr)

	require.NotNil(t, mockManager.Install(mockServer).Err)
}

func Test_Prune_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface
	mockInterface.On("Prune", mockServer, mockLocal, mockLogger, mockTimeout).Return(install.PruneErr{})
	require.Equal(t, install.PruneErr{}, mockManager.Prune(mockServer))
}

func Test_PruneErr_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface

	mockInterface.On("Prune", mockServer, mockLocal, mockLogger, mockTimeout).Return(mockPruneErr)
	require.NotNil(t, mockManager.Prune(mockServer))
}

func Test_Detect_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface
	mockInterface.On("Detect", "",
		mockServer, mockLocal, mockLogger, mockTimeout).Return(nil)

	require.Nil(t, mockManager.Detect(mockServer).Err)
}

func Test_HandPackage_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface
	mockInterface.On("HandPackage",
		mockServer, mockPackagePath, mockLocal, mockLogger, mockTimeout).
		Return(nil)

	require.Nil(t, mockManager.HandPackage(mockServer).Err)
}

func Test_HandPackageErr_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockManager.Handler = mockInterface
	mockInterface.On("HandPackage",
		mockServer, mockPackagePath, mockLocal, mockLogger, mockTimeout).
		Return(mockTransferErr)

	_, ok := mockManager.HandPackage(mockServer).Err.(install.TransferPackageErr)

	require.Equal(t, true, ok)
}

func Test_ParseYamlTypeErr(t *testing.T) {
	content := `
docker:
  package: 
    - docker-19.03.15.tgz   # 二进制安装包目录
  preserveDir: /data/lib/docker  # docker数据持久化目录
  insecureRegistries: # 非https仓库列表
    - gcr.azk8s.cn
    - quay.azk8s.cn
  registryMirrors:               # 镜像源
`
	mockManager.ConfigContent = []byte(content)
	re, err := mockManager.Parse()
	_, ok := err.(*yaml.TypeError)
	require.Equal(t, true, ok)
	require.Equal(t, 0, len(re.Servers))
}

func Test_ParseServerList(t *testing.T) {
	content := `
docker:
  package: docker-19.03.15.tgz   # 二进制安装包目录
  preserveDir: /data/lib/docker  # docker数据持久化目录
  insecureRegistries: # 非https仓库列表
    - gcr.azk8s.cn
    - quay.azk8s.cn
  registryMirrors:               # 镜像源
    - docker.a.io
    - docker.b.io
server:
  - host: 
     - 10.10.10.1-3
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
`
	mockManager.ConfigContent = []byte(content)
	re, err := mockManager.Parse()
	require.Equal(t, "docker-19.03.15.tgz", mockManager.Package)
	require.Equal(t, mockMirrorsSlice, mockManager.Mirrors)
	require.Equal(t, mockInsecureRegistriesSlice, mockManager.InsecureRegistries)
	require.Nil(t, err)
	require.Equal(t, 3, len(re.Servers))
}

// todo: yaml反序列化存在问题：反序列化[]byte，数组类型数据有影响
func Test_ParseServerListErr(t *testing.T) {
	content := `
docker:
  package: docker-19.03.15.tgz   # 二进制安装包目录
  preserveDir: /data/lib/docker  # docker数据持久化目录
  insecureRegistries: # 非https仓库列表
    - gcr.azk8s.cn
    - quay.azk8s.cn
  registryMirrors:               # 镜像源
    - docker.a.io
    - docker.b.io
server:
  - host: 
      - 10.10.10.1-3
  - username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
`
	mockManager.Servers = nil
	mockManager.ConfigContent = []byte(content)
	re, err := mockManager.Parse()
	_, ok := err.(install.ParseServerListErr)
	require.Equal(t, true, ok)
	require.Equal(t, 0, len(re.Servers))
}

func Test_getHandlerInterface(t *testing.T) {

	var h HandlerInterface
	r := getHandlerInterface(h)
	require.Equal(t, new(Handler), r)

	h2 := &mocks.HandlerInterface{}
	r2 := getHandlerInterface(h2)
	require.Equal(t, h2, r2)
}

func Test_InstallDocker_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	item := command.OperationItem{
		B:          []byte(mockInstallContent),
		Logger:     mockLogger,
		Interface:  mockInterface,
		Mock:       false,
		LocalRun:   false,
		SSHTimeout: mockTimeout,
	}

	mockInterface.On("Detect", "", mockServer, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Prune", mockServer, mockLocal, mockLogger, mockTimeout).Return(install.PruneErr{})
	mockInterface.On("HandPackage", mockServer, mockPackagePath, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Install", mockInstallShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.InstallErr{})
	mockInterface.On("SetUpRuntime", mockSetRuntimeShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.SetUpRuntimeErr{})
	mockInterface.On("SetConfig", mockSetConfigShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.SetConfigErr{})
	mockInterface.On("SetSystemd", mockSetSystemdShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.SetSystemdErr{})
	mockInterface.On("Boot", mockServer, mockLocal, mockLogger, mockTimeout).Return(install.BootErr{})

	runErr := Install(item)
	require.Nil(t, runErr.Err)
}

func Test_InstallDocker_BootErr_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	item := command.OperationItem{
		B:          []byte(mockInstallContent),
		Logger:     mockLogger,
		Interface:  mockInterface,
		Mock:       false,
		LocalRun:   false,
		SSHTimeout: mockTimeout,
	}

	mockInterface.On("Detect", "", mockServer, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Prune", mockServer, mockLocal, mockLogger, mockTimeout).Return(install.PruneErr{})
	mockInterface.On("HandPackage", mockServer, mockPackagePath, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Install", mockInstallShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.InstallErr{})
	mockInterface.On("SetUpRuntime", mockSetRuntimeShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.SetUpRuntimeErr{})
	mockInterface.On("SetConfig", mockSetConfigShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.SetConfigErr{})
	mockInterface.On("SetSystemd", mockSetSystemdShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.SetSystemdErr{})
	mockInterface.On("Boot", mockServer, mockLocal, mockLogger, mockTimeout).Return(mockBootErr)

	_, ok := Install(item).Err.(install.BootErr)
	require.Equal(t, true, ok)
}

func Test_InstallDocker_SetSystemdErr_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	item := command.OperationItem{
		B:          []byte(mockInstallContent),
		Logger:     mockLogger,
		Interface:  mockInterface,
		Mock:       false,
		LocalRun:   false,
		SSHTimeout: mockTimeout,
	}

	mockInterface.On("Detect", "", mockServer, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Prune", mockServer, mockLocal, mockLogger, mockTimeout).Return(install.PruneErr{})
	mockInterface.On("HandPackage", mockServer, mockPackagePath, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Install", mockInstallShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.InstallErr{})
	mockInterface.On("SetUpRuntime", mockSetRuntimeShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.SetUpRuntimeErr{})
	mockInterface.On("SetConfig", mockSetConfigShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.SetConfigErr{})
	mockInterface.On("SetSystemd", mockSetSystemdShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(mockSetSystemdErr)

	_, ok := Install(item).Err.(install.SetSystemdErr)
	require.Equal(t, true, ok)
}

func Test_InstallDocker_SetConfigErr_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	item := command.OperationItem{
		B:          []byte(mockInstallContent),
		Logger:     mockLogger,
		Interface:  mockInterface,
		Mock:       false,
		LocalRun:   false,
		SSHTimeout: mockTimeout,
	}

	mockInterface.On("Detect", "", mockServer, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Prune", mockServer, mockLocal, mockLogger, mockTimeout).Return(install.PruneErr{})
	mockInterface.On("HandPackage", mockServer, mockPackagePath, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Install", mockInstallShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.InstallErr{})
	mockInterface.On("SetUpRuntime", mockSetRuntimeShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.SetUpRuntimeErr{})
	mockInterface.On("SetConfig", mockSetConfigShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(mockSetConfigErr)

	_, ok := Install(item).Err.(install.SetConfigErr)
	require.Equal(t, true, ok)
}

func Test_InstallDocker_SetUpRuntimeErr_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	item := command.OperationItem{
		B:          []byte(mockInstallContent),
		Logger:     mockLogger,
		Interface:  mockInterface,
		Mock:       false,
		LocalRun:   false,
		SSHTimeout: mockTimeout,
	}

	mockInterface.On("Detect", "", mockServer, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Prune", mockServer, mockLocal, mockLogger, mockTimeout).Return(install.PruneErr{})
	mockInterface.On("HandPackage", mockServer, mockPackagePath, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Install", mockInstallShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(install.InstallErr{})
	mockInterface.On("SetUpRuntime", mockSetRuntimeShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(mockSetRuntimeErr)

	_, ok := Install(item).Err.(install.SetUpRuntimeErr)
	require.Equal(t, true, ok)
}

func Test_InstallDocker_InstallFuncErr_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	item := command.OperationItem{
		B:          []byte(mockInstallContent),
		Logger:     mockLogger,
		Interface:  mockInterface,
		Mock:       false,
		LocalRun:   false,
		SSHTimeout: mockTimeout,
	}

	mockInterface.On("Detect", "", mockServer, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Prune", mockServer, mockLocal, mockLogger, mockTimeout).Return(install.PruneErr{})
	mockInterface.On("HandPackage", mockServer, mockPackagePath, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Install", mockInstallShell, mockServer, mockLocal, mockLogger, mockTimeout).Return(mockInstallErr)

	_, ok := Install(item).Err.(install.InstallErr)
	require.Equal(t, true, ok)
}

func Test_InstallDocker_TransferPackageErr_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	item := command.OperationItem{
		B:          []byte(mockInstallContent),
		Logger:     mockLogger,
		Interface:  mockInterface,
		Mock:       false,
		LocalRun:   false,
		SSHTimeout: mockTimeout,
	}

	mockInterface.On("Detect", "", mockServer, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Prune", mockServer, mockLocal, mockLogger, mockTimeout).Return(install.PruneErr{})
	mockInterface.On("HandPackage", mockServer, mockPackagePath, mockLocal, mockLogger, mockTimeout).Return(mockTransferErr)

	_, ok := Install(item).Err.(install.TransferPackageErr)
	require.Equal(t, true, ok)
}

func Test_InstallDocker_PruneErr_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	item := command.OperationItem{
		B:          []byte(mockInstallContent),
		Logger:     mockLogger,
		Interface:  mockInterface,
		Mock:       false,
		LocalRun:   false,
		SSHTimeout: mockTimeout,
	}

	mockInterface.On("Detect", "", mockServer, mockLocal, mockLogger, mockTimeout).Return(nil)
	mockInterface.On("Prune", mockServer, mockLocal, mockLogger, mockTimeout).Return(mockPruneErr)

	_, ok := Install(item).Err.(install.PruneErr)
	require.Equal(t, true, ok)
}

func Test_InstallDocker_DetectErr_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	item := command.OperationItem{
		B:          []byte(mockInstallContent),
		Logger:     mockLogger,
		Interface:  mockInterface,
		Mock:       false,
		LocalRun:   false,
		SSHTimeout: mockTimeout,
	}

	mockInterface.On("Detect", "", mockServer, mockLocal, mockLogger, mockTimeout).Return(mockDetectErr)

	_, ok := Install(item).Err.(install.DetectErr)
	require.Equal(t, true, ok)
}

func Test_InstallDocker_ParseErr_Mock(t *testing.T) {

	content := `
server:
  - host: 1.1.1.1
    - username: root
    password: 123456
    port: 22
docker:
 package: docker.tgz   # 二进制安装包目录
 preserveDir: /data/docker  # docker数据持久化目录
 insecureRegistries: # 非https仓库列表
   - gcr.azk8s.cn
   - quay.azk8s.cn
 registryMirrors:               # 镜像源
   - docker.a.io
   - docker.b.io
`
	mockInterface := &mocks.HandlerInterface{}
	item := command.OperationItem{
		B:         []byte(content),
		Logger:    mockLogger,
		Interface: mockInterface,
		Mock:      false,
		LocalRun:  false,
	}

	err := Install(item).Err
	require.NotNil(t, err)
}

func Test_InstallDocker_Local(t *testing.T) {
	content := `
server:
  - host: 1.1.1.1
    username: root
    password: 123456
    port: 22
docker:
 package: docker.tgz   # 二进制安装包目录
 preserveDir: /data/docker  # docker数据持久化目录
 insecureRegistries: # 非https仓库列表
   - gcr.azk8s.cn
   - quay.azk8s.cn
 registryMirrors:               # 镜像源
   - docker.a.io
   - docker.b.io
`
	//mockInterface := &mocks.HandlerInterface{}
	item := command.OperationItem{
		B:      []byte(content),
		Logger: mockLogger,
		//Interface: mockInterface,
		Mock:     false,
		LocalRun: true,
	}

	err := Install(item).Err
	require.NotNil(t, err)
}
