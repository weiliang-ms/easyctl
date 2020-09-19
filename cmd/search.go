package cmd

import (
	"easyctl/sys"
	"easyctl/util"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	searchCmd.AddCommand(searchPortCmd)
	rootCmd.AddCommand(searchCmd)
}

// set 命令合法参数
var searchValidArgs = []string{"port"}

// 查询命令search指令集
var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "search something through easyctl",
	Long:    `Search port is lientened...`,
	Example: "search port 8888",
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("配置功能...")
	},
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: searchValidArgs,
}

// 查询端口监听指令
var searchPortCmd = &cobra.Command{
	Use:     "port [value]",
	Short:   "search port listen status",
	Long:    `Search port is lientened...`,
	Example: "search port 8888",
	Run: func(cmd *cobra.Command, args []string) {
		searchPortStatus(args)
	},
	Args: cobra.MinimumNArgs(1),
}

// 查询端口监听状态
func searchPortStatus(args []string) {

	util.PrintSuccessfulMsg("#### 查询端口：%s监听状态 ####")
	err := sys.SearchPortStatus(args[0])

	if err != nil {
		util.PrintFailureMsg(fmt.Sprintf("\n[notOnListening] 端口：%s未被监听或您输入的端口格式不对，请检查是否为1~65535\n", args[1]))
	} else {
		util.PrintSuccessfulMsg(fmt.Sprintf("\n[onListening] 端口：%s处于被监听状态\n", args[1]))
	}
}
