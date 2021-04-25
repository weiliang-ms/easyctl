package close

import (
	"github.com/spf13/cobra"
)

var (
	forever        bool
	serverListFile string
)

func init() {
	RootCmd.PersistentFlags().BoolVarP(&forever, "forever", "", true, "是否永久关闭服务")
	RootCmd.PersistentFlags().StringVarP(&serverListFile, "server-list", "", "server.yaml", "服务器列表连接信息")
	RootCmd.AddCommand(closeFirewallCmd)
	RootCmd.AddCommand(closeSeLinuxCmd)
}

// close命令
var RootCmd = &cobra.Command{
	Use:     "close [OPTIONS] [flags]",
	Example: "\neasyctl close firewalld --forever",
	Run: func(cmd *cobra.Command, args []string) {
	},
	ValidArgs: []string{"firewall", "selinux"},
	Args:      cobra.ExactValidArgs(1),
}
