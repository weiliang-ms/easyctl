package cmd

import (
	"easycfg/sys"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

const dns = "dns"
const yum = "yum源"

// todo
var setHelpContent = "这是设置的帮助内容..."
var configFailed = "配置失败..."
var configSuccess = "配置成功..."
var successBanner = "[success]"
var failedBanner = "[failed]"

var missingParameterErr = errors.New("缺少参数...")

func init() {
	rootCmd.AddCommand(setCmd)
}

// set 命令合法参数
var setValidArgs = []string{"dns", "yum", "hostname"}

// 输出easycfg版本
var setCmd = &cobra.Command{
	Use:     "set",
	Short:   "set something through easycfg",
	Long:    `Set DNS address,hostname,yum address...`,
	Example: "set dns 114.114.114.114",
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("配置功能...")
		analyseArgs(args)
	},
	ValidArgs: setValidArgs,
}

func analyseArgs(args []string) {
	if len(args) == 0 {
		setHelp()
	} else {
		switch args[0] {
		case "dns":
			setDNS(args)
		case "yum":
			setYUM(args)
		case "hostname":
			setHostname(args)
		default:
			setHelp()
		}
	}
}

// 配置dns
func setDNS(args []string) {
	if len(args) < 2 {
		fmt.Println(missingParameterErr)
	} else {
		fmt.Printf("#### 配置dns地址：%s ####\n", args[1])
		err, _ := sys.SetDNS(args[1])
		if err != nil {
			fmt.Println(failedBanner + " " + dns + configFailed + ": " + err.Error())
		} else {
			fmt.Println(successBanner + dns + configSuccess)
		}
	}
}

// 配置yum
func setYUM(args []string) {
	if len(args) < 2 {
		sys.SetAliYUM()
	} else {
		switch args[1] {
		case "ali":
			sys.SetAliYUM()
		default:
			// todo
			fmt.Println("暂不支持该mirror...")
		}
	}
}

// 配置hostname

func setHostname(args []string) {
	if len(args) < 2 {
		// todo
		fmt.Println("设置hostname帮助逻辑...")
	} else {
		sys.SetHostname(args[1])
	}
}
func setHelp() {
	fmt.Println(setHelpContent)
}
