package set

import (
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/set"
)

// PasswordLessCmd 主机互信
var passwordLessCmd = &cobra.Command{
	Use:     "password-less [flags]",
	Short:   "配置主机互信",
	Example: "\neasyctl set password-less --server-list=server.yaml",
	Args:    cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := Set(Entity{Cmd: cmd, Fnc: set.PasswordLess}); runErr != nil {
			panic(runErr)
		}
	},
}
