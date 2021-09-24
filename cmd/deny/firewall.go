package deny

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/deny"
)

// 关闭防火墙
var denyFirewallCmd = &cobra.Command{
	Use:   "firewall [flags]",
	Short: "关闭防火墙",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := Deny(Entity{Cmd: cmd, Fnc: deny.Firewall}); runErr != nil {
			panic(runErr)
		}
	},
}
