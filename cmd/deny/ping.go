package deny

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/deny"
)

// 关闭 ping
var denyPingCmd = &cobra.Command{
	Use:   "ping [flags]",
	Short: "禁ping",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := Deny(Entity{Cmd: cmd, Fnc: deny.Ping}); runErr != nil {
			panic(runErr)
		}
	},
}
