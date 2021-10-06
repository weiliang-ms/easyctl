package set

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
	Use:   "set [OPTIONS] [flags]",
	Short: "设置指令集",
	Args:  cobra.ExactValidArgs(1),
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	subCommands := []*cobra.Command{
		dnsCmd,
		hostResolveCmd,
		passwordLessCmd,
		newPasswordCmd,
		timeZoneCmd,
		ulimitCmd,
		yumRepoCmd,
	}

	for _, v := range subCommands {
		RootCmd.AddCommand(v)
	}
}