package clean

import "github.com/spf13/cobra"

var (
	serverListFile string
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&serverListFile, "server-list", "", "", "服务器列表连接信息")
	RootCmd.AddCommand(dnsCmd)
}

// clean命令
var RootCmd = &cobra.Command{
	Use:     "clean [OPTIONS] [flags]",
	Short:   "清理指令集",
	Example: "\neasyctl clean dns",
	Run: func(cmd *cobra.Command, args []string) {
	},
	ValidArgs: []string{"dns"},
	Args:      cobra.ExactValidArgs(1),
}
