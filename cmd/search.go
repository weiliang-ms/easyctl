package cmd

import (
	"easyctl/sys"
	"easyctl/util"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

// set 命令合法参数
var searchValidArgs = []string{"port"}

// 输出easyctl版本
var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "search something through easyctl",
	Long:    `Search port is lientened...`,
	Example: "search port 8888",
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("配置功能...")
		search(args)
	},
	//ValidArgs: setValidArgs,
}

func search(args []string) {
	if len(args) == 0 {
		// todo
		fmt.Println("search语句帮助信息")
	} else {
		switch args[0] {
		case "port", "p":
			searchPortStatus(args)
		default:
			fmt.Println("search语句帮助信息")
		}
	}
}

func searchPortStatus(args []string) {
	if len(args) < 2 {
		fmt.Println(missingParameterErr)
	} else {
		util.PrintSuccessfulMsg("#### 查询端口：%s监听状态 ####")
		err := sys.SearchPortStatus(args[1])
		if err != nil {
			util.PrintFailureMsg(fmt.Sprintf("\n[notOnListening] 端口：%s未被监听或您输入的端口格式不对，请检查是否为1~65535\n", args[1]))
		} else {
			util.PrintSuccessfulMsg(fmt.Sprintf("\n[onListening] 端口：%s处于被监听状态\n", args[1]))
		}
	}
}
