package install

// Interface install操作类型接口
type Interface interface {
	Detect() error
	Prune() error
	HandPackage() error
	Compile() error
}

type task func() error

// Install 安装指令通用函数
func install(i Interface) error {

	jobs := []task{
		i.Detect,
		i.Prune,
		i.HandPackage,
		i.Compile,
	}

	for _, v := range jobs {
		if err := v(); err != nil {
			return err
		}
	}

	return nil
}
