package scan

import (
	"github.com/spf13/cobra"
)

var (
	configFile string
)

// RootCmd 扫描命令
var RootCmd = &cobra.Command{
	Use:   "scan [flags]",
	Short: "扫描命令指令集",
	Args:  cobra.ExactValidArgs(1),
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(osCmd)
}
