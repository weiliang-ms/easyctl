package track

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
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(tailLogCmd)
}

// RootCmd track命令
var RootCmd = &cobra.Command{
	Use:   "track [flags]",
	Short: "追踪命令指令集",
	Args:  cobra.ExactValidArgs(1),
}
