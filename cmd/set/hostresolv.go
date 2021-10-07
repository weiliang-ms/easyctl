package set

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/set"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

// host解析
// 主机互信
var hostResolveCmd = &cobra.Command{
	Use:   "host-resolv [flags]",
	Short: "配置host解析",
	Args:  cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if err := command.SetExecutorDefault(command.Item{
			Cmd:            cmd,
			Fnc:            set.HostResolve,
			ConfigFilePath: configFile,
		}); err != nil {
			panic(err)
		}
	},
}
