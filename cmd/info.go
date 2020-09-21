package cmd

import (
	"easyctl/sys"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	infoCmd.AddCommand(osInfoCmd)
	infoCmd.AddCommand(sysInfoCmd)
	rootCmd.AddCommand(infoCmd)
}

// set 命令合法参数
var infoValidArgs = []string{"os", "system"}

// 查询命令search指令集
var infoCmd = &cobra.Command{
	Use:     "info",
	Short:   "print information about current system through easyctl",
	Long:    `Search port is lientened...`,
	Example: "info os",
	Run: func(cmd *cobra.Command, args []string) {

	},
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: searchValidArgs,
}

// 操作系统版本信息
var osInfoCmd = &cobra.Command{
	Use:     "os",
	Short:   "print system version",
	Long:    `print system version...`,
	Example: "info os",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("当前系统版本：%+v\n", sys.SystemInfoObject.OSVersion.ReleaseContent)
	},
}

// 操作系统整体信息
var sysInfoCmd = &cobra.Command{
	Use:     "system",
	Short:   "print system version",
	Long:    `print system version...`,
	Example: "info os",
	Aliases: []string{"sys"},
	Run: func(cmd *cobra.Command, args []string) {
		sys.PrintSystemInfo()
	},
}
