package close

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/close"
)

// 关闭 ping
var closePingCmd = &cobra.Command{
	Use:   "ping [flags]",
	Short: "禁ping",
	Run: func(cmd *cobra.Command, args []string) {
		denyPing()
	},
}

func denyPing() {
	close := close.Actuator{
		ServerListFile: serverListFile,
	}
	close.Ping()
}
