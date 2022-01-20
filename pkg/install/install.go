package install

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

type BootErr struct {
	Host string
	Err  error
}

type SetSystemdErr struct {
	Host string
	Err  error
}

type SetConfigErr struct {
	Host string
	Err  error
}

type SetUpRuntimeErr struct {
	Host string
	Err  error
}

type InstallErr struct {
	Host string
	Err  error
}

type DetectErr struct {
	Host string
	Err  error
}

type PruneErr struct {
	Host string
	Err  error
}

type TransferPackageErr struct {
	Host string
	Path string
	Err  error
}

type ParseServerListErr struct {
	Err error
}

// Interface install操作类型接口
type Interface interface {
	Parse() command.RunErr
	SetValue() command.RunErr
	Detect() command.RunErr
	Prune() command.RunErr
	HandPackage() command.RunErr
	Install() command.RunErr
	SetUpRuntime() command.RunErr
	Config() command.RunErr
	SetService() command.RunErr // 开机自启动
	Boot() command.RunErr
	CloseFirewall() command.RunErr
	Init() command.RunErr
	Print() command.RunErr // 输出安装信息
}

type task func() command.RunErr

// Install 安装指令通用函数
func install(i Interface) command.RunErr {

	jobs := []task{
		i.Parse,
		i.SetValue,
		i.Detect,
		i.Prune,
		i.HandPackage,
		i.Install,
		i.SetUpRuntime,
		i.Config,
		i.SetService,
		i.Boot,
		i.CloseFirewall,
		i.Init,
		i.Print,
	}

	for _, v := range jobs {
		if runErr := v(); runErr.Err != nil {
			return runErr
		}
	}

	return command.RunErr{}
}

func (bootErr BootErr) Error() string {
	return fmt.Sprintf("[%s] 启动失败 -> %s", bootErr.Host, bootErr.Err)
}

func (setSystemdErr SetSystemdErr) Error() string {
	return fmt.Sprintf("[%s] 配置systemd失败 -> %s", setSystemdErr.Host, setSystemdErr.Err)
}

func (setConfigErr SetConfigErr) Error() string {
	return fmt.Sprintf("[%s] 配置docker失败 -> %s", setConfigErr.Host, setConfigErr.Err)
}

func (setUpRuntimeErr SetUpRuntimeErr) Error() string {
	return fmt.Sprintf("[%s] 配置docker运行时失败 -> %s", setUpRuntimeErr.Host, setUpRuntimeErr.Err)
}

func (installErr InstallErr) Error() string {
	return fmt.Sprintf("[%s] 安装docker失败 -> %s", installErr.Host, installErr.Err)
}

func (pruneErr PruneErr) Error() string {
	return fmt.Sprintf("[%s] 清理docker失败 -> %s", pruneErr.Host, pruneErr.Err)
}

func (detectErr DetectErr) Error() string {
	return fmt.Sprintf("[%s] 检测docker安装依赖失败 -> %s", detectErr.Host, detectErr.Err)
}

func (parseServerListErr ParseServerListErr) Error() string {
	return fmt.Sprintf("反序列化主机列表失败 -> %s", parseServerListErr.Err)
}

func (transferPackageErr TransferPackageErr) Error() string {
	return fmt.Sprintf("传输%s:/tmp/%s失败 -> %s",
		transferPackageErr.Host,
		transferPackageErr.Path,
		transferPackageErr.Err)
}
