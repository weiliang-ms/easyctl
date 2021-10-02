package upgrade

import "github.com/spf13/cobra"

var (
	offline        bool
	kernelVersion  string
	filePath       string
	serverListFile string
)

func init() {
	Cmd.PersistentFlags().BoolVarP(&offline, "offline", "", false, "是否离线安装")
	Cmd.PersistentFlags().StringVarP(&serverListFile, "server-list", "", "", "服务器列表连接信息")
	Cmd.AddCommand(upgradeKernelCmd)
	Cmd.AddCommand(upgradeOpensslCmd)
	Cmd.AddCommand(upgradeOpenSSHCmd)
}

// Cmd upgrade 命令
var Cmd = &cobra.Command{
	Use:     "upgrade [OPTIONS] [flags]",
	Short:   "更新指令集",
	Example: "\neasyctl upgrade kernel --kernel-version=lt\n",
	Args:    cobra.MinimumNArgs(1),
}
