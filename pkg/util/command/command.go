package command

import (
	// embed
	_ "embed"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
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
func SetExecutorDefault(entity ExecutorEntity, configFile string) error {

	if entity.DefaultConfig == nil {
		entity.DefaultConfig = config
	}

	if configFile == "" {
		logrus.Infof("检测到配置文件参数为空，生成配置文件样例 -> %s", constant.ConfigFile)
		_ = os.WriteFile(constant.ConfigFile, entity.DefaultConfig, 0666)
		os.Exit(0)
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
