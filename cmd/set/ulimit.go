package set

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/set"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

// 文件描述符
var ulimitCmd = &cobra.Command{
	Use:   "ulimit [flags]",
	Short: "配置ulimit",
	Args:  cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if err := command.SetExecutorDefault(command.Item{
			Cmd:            cmd,
			Fnc:            set.Ulimit,
			ConfigFilePath: configFile,
		}); err.Err != nil {
			panic(err.Err)
		}
	},
}
