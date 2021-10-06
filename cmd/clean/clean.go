package clean

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

// RootCmd set命令
var RootCmd = &cobra.Command{
	Use:   "clean [OPTIONS] [flags]",
	Short: "清理指令集",
	Args:  cobra.ExactValidArgs(1),
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(cleanDnsCmd)
}
