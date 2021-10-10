package install

import "github.com/weiliang-ms/easyctl/pkg/util/command"

// Interface install操作类型接口
type Interface interface {
	Parse() command.RunErr
	SetValue() command.RunErr
	Detect() command.RunErr
	Prune() command.RunErr
	HandPackage() command.RunErr
	Compile() command.RunErr
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
		i.Compile,
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
