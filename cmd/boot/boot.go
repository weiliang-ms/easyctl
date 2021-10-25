package boot

import (
	"github.com/spf13/cobra"
)

var (
	configFile string
)

// RootCmd boot命令
var RootCmd = &cobra.Command{
	Use:   "boot [OPTIONS] [flags]",
	Short: "启动指令集",
	Args:  cobra.ExactValidArgs(1),
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(appWithCGroupsCmd)
}
