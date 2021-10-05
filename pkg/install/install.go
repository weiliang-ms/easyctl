package install

// Interface install操作类型接口
type Interface interface {
	Parse() error
	Detect() error
	Prune() error
	HandPackage() error
	Compile() error
	SetUpRuntime() error
	Config() error
	Boot() error
	CloseFirewall() error
}

type task func() error

// Install 安装指令通用函数
func install(i Interface) error {

	jobs := []task{
		i.Parse,
		i.Detect,
		i.Prune,
		i.HandPackage,
		i.Compile,
		i.SetUpRuntime,
		i.Config,
		i.Boot,
		i.CloseFirewall,
	}

	for _, v := range jobs {
		if err := v(); err != nil {
			return err
		}
	}

	return nil
}
