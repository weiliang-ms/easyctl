package deny

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/deny"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

// 关闭 ping
var denyPingCmd = &cobra.Command{
	Use:   "ping [flags]",
	Short: "禁ping",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(command.Item{Cmd: cmd, Fnc: deny.Ping}); runErr != nil {
			panic(runErr)
		}
	},
}
