package harden

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

//go:embed asset/config.yaml
var config []byte

// RootCmd 安全加固命令
var RootCmd = &cobra.Command{
	Use:   "harden [OPTIONS] [flags]",
	Short: "安全加固指令",
	Args:  cobra.ExactValidArgs(1),
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(osHardenCmd)
}
