package exec

import (
	// embed
	_ "embed"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/runner"
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
	RootCmd.Flags().StringVarP(&configFile, "config", "c", "", "配置文件")
}

// RootCmd close命令
var RootCmd = &cobra.Command{
	Use:     "exec [flags]",
	Short:   "执行命令指令集",
	Example: "\neasyctl exec -c config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := Exec(Entity{Cmd: cmd, Fnc: runner.Run}); runErr != nil {
			panic(runErr)
		}
	},
}

// Entity 实体
type Entity struct {
	Cmd           *cobra.Command
	Fnc           func(b []byte, logger *logrus.Logger) error
	DefaultConfig []byte
}

// Exec 组装执行器
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
