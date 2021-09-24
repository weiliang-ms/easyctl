package set

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/set"
)

// host解析
// 主机互信
var hostResolveCmd = &cobra.Command{
	Use:     "host-resolv [flags]",
	Short:   "配置host解析",
	Example: "\neasyctl set host-resolv --server-list=server.yaml",
	Args:    cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := Set(Entity{Cmd: cmd, Fnc: set.HostResolve}); runErr != nil {
			panic(runErr)
		}
	},
}
