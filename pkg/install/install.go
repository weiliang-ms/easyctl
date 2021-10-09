package install

// Interface install操作类型接口
type Interface interface {
	Parse() error
	SetValue() error
	Detect() error
	Prune() error
	HandPackage() error
	Compile() error
	SetUpRuntime() error
	Config() error
	SetService() error // 开机自启动
	Boot() error
	CloseFirewall() error
	Init() error
	Print() error // 输出安装信息
}

type task func() error

// Install 安装指令通用函数
func install(i Interface) error {

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
		if err := v(); err != nil {
			return err
		}
	}

	return nil
}
