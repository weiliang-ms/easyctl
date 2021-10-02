package track

import (
	// embed
	_ "embed"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"k8s.io/klog"
	"os"
)

var (
	configFile string
)

//go:embed asset/executor.yaml
var config []byte

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(tailLogCmd)
}

// RootCmd track命令
var RootCmd = &cobra.Command{
	Use:   "track [flags]",
	Short: "追踪命令指令集",
	Args:  cobra.ExactValidArgs(1),
}

// Entity track指令实体
type Entity struct {
	Cmd           *cobra.Command
	Fnc           func(b []byte, logger *logrus.Logger) error
	DefaultConfig []byte
}

// Exec track执行实体
func Exec(entity Entity) error {

	if entity.DefaultConfig == nil {
		entity.DefaultConfig = config
	}

	if configFile == "" {
		klog.Infof("检测到配置文件参数为空，生成配置文件样例 -> %s", util.ConfigFile)
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
		klog.Fatalf("读取配置文件失败")
	}

	logger := logrus.New()
	if debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return entity.Fnc(b, logger)
}
