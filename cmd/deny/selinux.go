package deny

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/deny"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

// 关闭selinux
var denySelinuxCmd = &cobra.Command{
	Use:   "selinux [flags]",
	Short: "关闭selinux",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(command.Item{Cmd: cmd, Fnc: deny.Selinux}); runErr != nil {
			panic(runErr)
		}
	},
}
