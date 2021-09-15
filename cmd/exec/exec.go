package exec

import (
	_ "embed"
	"github.com/spf13/cobra"
	exec2 "github.com/weiliang-ms/easyctl/pkg/exec"
	"k8s.io/klog"
)

var (
	configFile string
)

func init() {
	RootCmd.Flags().StringVarP(&configFile, "config", "c", "", "配置文件")
}

// RootCmd close命令
var RootCmd = &cobra.Command{
	Use:     "exec [flags]",
	Short:   "执行命令指令集",
	Example: "\neasyctl exec -c config.yaml",
	Run: func(cmd *cobra.Command, args []string) {

		if err := exec(); err != nil {
			klog.Fatalln(err)
		}
	},
}

func exec() error {
	return exec2.Run(configFile, 1)
}
