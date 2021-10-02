package add

import (
	// embed
	_ "embed"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"os"
)

var (
	configFile string
)

type Entity struct {
	Cmd           *cobra.Command
	Fnc           func(b []byte, logger *logrus.Logger) error
	DefaultConfig []byte
}

//go:embed asset/config.yaml
var config []byte

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(addUserCmd)
}

// RootCmd add命令
var RootCmd = &cobra.Command{
	Use:   "add [flags]",
	Short: "添加指令集",
	Args:  cobra.ExactValidArgs(1),
}

// Add 组装执行器
func Add(entity Entity) error {

	if configFile == "" {
		logrus.Infof("检测到配置文件参数为空，生成配置文件样例 -> %s", util.ConfigFile)
		_ = os.WriteFile(util.ConfigFile, entity.DefaultConfig, 0666)
		os.Exit(0)
	}

	flagset := entity.Cmd.Parent().Parent().PersistentFlags()
	debug, err := flagset.GetBool("debug")
	if err != nil {
		fmt.Println(err)
	}

	b, readErr := os.ReadFile(configFile)
	if readErr != nil {
		logrus.Fatalf("读取配置文件失败")
	}

	logger := logrus.New()
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	return entity.Fnc(b, logger)
}

func (entity *Entity) setDefault() {
	if entity.DefaultConfig == nil {
		entity.DefaultConfig = config
	}
}
