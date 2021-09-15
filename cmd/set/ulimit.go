package set

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/set/ulimit"
)

// 文件描述符
var ulimitCmd = &cobra.Command{
	Use:     "ulimit [flags]",
	Short:   "配置ulimit",
	Example: "\neasyctl set ulimit -c config.yaml",
	Args:    cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		err := ulimit.Config(configFile, level())
		if err != nil {
			panic(err)
		}
	},
}
