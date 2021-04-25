package close

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/close"
)

// 关闭防火墙
var closeFirewallCmd = &cobra.Command{
	Use:   "firewall [flags]",
	Short: "easyctl close firewall [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		c := close.Closer{
			ServerListFilePath: serverListFile,
			Forever:            forever,
		}
		c.Firewall()
	},
}
