package export

import (
	"github.com/spf13/cobra"
)

var (
	packageName    string
	serverListFile string
	configFilePath string
)

// RootCmd export命令
var RootCmd = &cobra.Command{
	Use:   "export [OPTIONS] [flags]",
	Short: "导出指令集",
	Run: func(cmd *cobra.Command, args []string) {
	},
	ValidArgs: []string{"yum-repo"},
	Args:      cobra.ExactValidArgs(1),
}

func init() {
	// todo: 变更入参 --server-list -> --config or -c
	RootCmd.PersistentFlags().StringVarP(&serverListFile, "server-list", "", "", "服务器批量连接信息")
	RootCmd.AddCommand(imageCmd)
	RootCmd.AddCommand(localImagesCmd)
	RootCmd.AddCommand(chartCmd)
}
