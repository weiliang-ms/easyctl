package add

import (
	"github.com/spf13/cobra"
)

var (
	Nologin        bool
	username       string
	password       string
	serverListFile string
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&serverListFile, "server-list", "", "", "服务器列表连接信息")
	RootCmd.AddCommand(addUserCmd)
}

// RootCmd add命令
var RootCmd = &cobra.Command{
	Use:     "add [OPTIONS] [flags]",
	Short:   "添加指令集",
	Example: "\neasyctl add user user1 password",
	Run: func(cmd *cobra.Command, args []string) {
	},
	ValidArgs: []string{"user"},
	Args:      cobra.ExactValidArgs(1),
}
