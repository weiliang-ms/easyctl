package exec

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
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

	},
}
