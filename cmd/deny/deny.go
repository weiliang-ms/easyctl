package deny

import (
	// embed
	_ "embed"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

//go:embed asset/config.yaml
var config []byte

// Entity 实体
type Entity struct {
	Cmd *cobra.Command
	Fnc func(b []byte, logger *logrus.Logger) error
}

// RootCmd 禁用命令
var RootCmd = &cobra.Command{
	Use:   "deny [OPTIONS] [flags]",
	Short: "禁用指令集",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(denyFirewallCmd)
	RootCmd.AddCommand(denySelinuxCmd)
	RootCmd.AddCommand(denyPingCmd)
}
