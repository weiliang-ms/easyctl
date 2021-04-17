package close

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/runner"
)

func init() {
	closeFirewallShell = "systemctl disable firewalld --now"
	disableFirewallShell = "systemctl stop firewalld"
	closeFirewallCmd.Flags().BoolVarP(&forever, "forever", "", true, "是否永久关闭防火墙")
	closeFirewallCmd.Flags().BoolVarP(&remote, "remote", "", false, "是否关闭远程主机防火墙")
	closeFirewallCmd.Flags().StringVarP(&serverListFile, "server-list", "", "server.yaml", "服务器列表连接信息")

}

var (
	closeFirewallShell   string
	disableFirewallShell string
)

// 关闭防火墙
var closeFirewallCmd = &cobra.Command{
	Use:     "firewall [flags]",
	Short:   "easyctl close firewall [flags]",
	Example: "\neasyctl close firewall --forever",
	Run: func(cmd *cobra.Command, args []string) {
		closeFirewall()
	},
}

// 关闭防火墙
func closeFirewall() {
	var list []runner.Server
	if remote {
		list = runner.ParseServerList(serverListFile).Server
	}
	close(cmd(), list)
}

// todo: 优化if else
func cmd() string {
	if forever {
		return disableFirewallShell
	}
	return closeFirewallShell
}
