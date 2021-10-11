package deny

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/deny"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"log"
)

// 关闭防火墙
var denyFirewallCmd = &cobra.Command{
	Use:   "firewall [flags]",
	Short: "关闭防火墙",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(command.Item{
			Cmd: cmd, Fnc: deny.Firewall}); runErr.Err != nil {
			log.Println(runErr.Msg)
			panic(runErr.Err)
		}
	},
}
