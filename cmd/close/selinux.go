package close

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/close"
)

// 关闭selinux
var closeSeLinuxCmd = &cobra.Command{
	Use:   "selinux [flags]",
	Short: "关闭selinux",
	Run: func(cmd *cobra.Command, args []string) {
		c := close.Actuator{
			ServerListFile: serverListFile,
		}
		c.SeLinux()
	},
}
