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

// ExecutorEntity executor实体
type ExecutorEntity struct {
	Cmd           *cobra.Command
	Fnc           func(b []byte, logger *logrus.Logger) error
	DefaultConfig []byte
}

// SetExecutorDefault executor赋值
func SetExecutorDefault(entity ExecutorEntity, configFile string) (err error) {

	callerName := "github.com/weiliang-ms/easyctl/pkg/util/command.TestSetExecutorDefaultReturnErr"
	defer errors.IgnoreErrorFromCaller(2, callerName, &err)

	if entity.DefaultConfig == nil {
		entity.DefaultConfig = config
	}

	if configFile == "" {
		logrus.Infof("生成配置文件样例, 请携带 -c 参数重新执行 -> %s", constant.ConfigFile)
		_ = os.WriteFile(constant.ConfigFile, entity.DefaultConfig, 0666)
		return nil
	}

	flagset := entity.Cmd.Parent().Parent().PersistentFlags()
	debug, err := flagset.GetBool("debug")
	if err != nil {
		return err
	}

	b, readErr := os.ReadFile(configFile)
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
	return entity.Fnc(b, logger)
}
