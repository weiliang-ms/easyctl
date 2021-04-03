package set

import (
	"easyctl/sys"
	"github.com/spf13/cobra"
)

// 配置hostname子命令
var setHostnameCmd = &cobra.Command{
	Use:     "hostname [value]",
	Short:   "easyctl set hostname [value]",
	Example: "\neasyctl set hostname nginx-server1",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		setHostname()
	},
}

// 配置hostname
func setHostname() {
	sys.SetHostname(value)
}
