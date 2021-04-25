package close

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/close"
)

// 关闭selinux
var closeSeLinuxCmd = &cobra.Command{
	Use:   "selinux [flags]",
	Short: "easyctl close selinux",
	Run: func(cmd *cobra.Command, args []string) {
		c := close.Closer{
			ServerListFilePath: serverListFile,
			Forever:            forever,
		}
		c.SeLinux()
	},
}
