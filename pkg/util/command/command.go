package command

import (
	// embed
	_ "embed"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"github.com/weiliang-ms/easyctl/pkg/util/errors"
	"github.com/weiliang-ms/easyctl/pkg/util/log"
	"os"
)

//go:embed config.yaml
var config []byte

// Item 初始化赋值实体
type Item struct {
	Cmd            *cobra.Command
	DefaultConfig  []byte
	Fnc            func(item OperationItem) error
	ConfigFilePath string
	OptionFunc     map[string]interface{}
}

// OperationItem 操作实体
type OperationItem struct {
	B          []byte
	Logger     *logrus.Logger
	OptionFunc map[string]interface{}
}

// SetExecutorDefault executor赋值
func SetExecutorDefault(item Item) (err error) {

	callerName := "github.com/weiliang-ms/easyctl/pkg/util/command.TestSetExecutorDefaultReturnErr"
	defer errors.IgnoreErrorFromCaller(2, callerName, &err)

	if item.DefaultConfig == nil {
		item.DefaultConfig = config
	}

	if item.ConfigFilePath == "" {
		logrus.Infof("生成配置文件样例, 请携带 -c 参数重新执行 -> %s", constant.ConfigFile)
		_ = os.WriteFile(constant.ConfigFile, item.DefaultConfig, 0666)
		return nil
	}

	flagset := item.Cmd.Parent().Parent().PersistentFlags()
	debug, err := flagset.GetBool("debug")
	if err != nil {
		return err
	}

	b, readErr := os.ReadFile(item.ConfigFilePath)
	if readErr != nil {
		return readErr
	}

	logger := logrus.New()
	if debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.SetFormatter(&log.CustomFormatter{})

	return item.Fnc(OperationItem{
		B:          b,
		Logger:     logger,
		OptionFunc: item.OptionFunc,
	})
}
