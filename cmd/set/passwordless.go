package set

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/set"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

// PasswordLessCmd 主机互信
var passwordLessCmd = &cobra.Command{
	Use:     "password-less [flags]",
	Short:   "配置主机互信",
	Example: "\neasyctl set password-less --server-list=server.yaml",
	Args:    cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if err := command.SetExecutorDefault(command.Item{
			Cmd:            cmd,
			Fnc:            set.PasswordLess,
			ConfigFilePath: configFile,
		}); err.Err != nil {
			panic(err.Err)
		}
	},
}
