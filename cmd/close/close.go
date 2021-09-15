package close

import (
	"github.com/spf13/cobra"
)

var (
	forever        bool
	serverListFile string
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&serverListFile, "server-list", "", "", "服务器列表连接信息")
	RootCmd.AddCommand(closeFirewallCmd)
	RootCmd.AddCommand(closeSeLinuxCmd)
	RootCmd.AddCommand(closePingCmd)
}

// RootCmd close命令
var RootCmd = &cobra.Command{
	Use:     "close [OPTIONS] [flags]",
	Short:   "关闭指令集",
	Example: "\neasyctl close firewalld --forever",
	Run: func(cmd *cobra.Command, args []string) {
	},
	ValidArgs: []string{"firewall", "selinux"},
	Args:      cobra.ExactValidArgs(1),
}
