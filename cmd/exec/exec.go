package exec

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

// RootCmd close命令
var RootCmd = &cobra.Command{
	Use:     "exec [flags]",
	Short:   "执行命令指令集",
	Example: "\neasyctl exec -c config.yaml",
	Args:    cobra.ExactValidArgs(1),
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(shellCmd)
	RootCmd.AddCommand(scpCmd)
	RootCmd.AddCommand(pingCmd)
}
