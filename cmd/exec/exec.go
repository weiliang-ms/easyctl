package exec

import (
	_ "embed"
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

		if configFile == "" {
			klog.Infof("检测到配置文件为空，生成配置文件样例 -> %s", util.ConfigFile)
			os.WriteFile(util.ConfigFile, config, 0666)
		}

		b, err := os.ReadFile(configFile)
		if err != nil {
			panic(err)
		}

		flagset := cmd.Parent().Parent().PersistentFlags()
		debug, err := flagset.GetBool("debug")
		if err != nil {
			panic(err)
		}

		if err := runner.Run(b, debug); err != nil {
			klog.Fatalln(err)
		}
	},
}
