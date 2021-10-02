package export

import (
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

// RootCmd export命令
var RootCmd = &cobra.Command{
	Use:   "export [OPTIONS] [flags]",
	Short: "导出指令集",
	Args:  cobra.ExactValidArgs(1),
}

func init() {
	// todo: 变更入参 --server-list -> --config or -c
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(imageCmd)
	RootCmd.AddCommand(localImagesCmd)
	RootCmd.AddCommand(chartCmd)
}

// Entity 执行指令实体
type Entity struct {
	Cmd           *cobra.Command
	Fnc           func(b []byte, logger *logrus.Logger) error
	DefaultConfig []byte
}

// Export 组装执行器
func Export(entity Entity) error {

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
