package set

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/set/passwordless"
)

// 主机互信
var passwordLessCmd = &cobra.Command{
	Use:     "password-less [flags]",
	Short:   "",
	Example: "\neasyctl set password-less --server-list=server.yaml",
	Args:    cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		passwordless.Config(configFile, level())
	},
}
