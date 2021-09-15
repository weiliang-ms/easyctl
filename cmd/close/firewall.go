package close

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/close"
)

// 关闭防火墙
var closeFirewallCmd = &cobra.Command{
	Use:   "firewall [flags]",
	Short: "关闭防火墙",
	Run: func(cmd *cobra.Command, args []string) {
		c := close.Actuator{
			ServerListFile: serverListFile,
		}
		c.Firewall()
	},
}
