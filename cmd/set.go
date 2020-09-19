package cmd

import (
	"easyctl/sys"
	"errors"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
)

const (
	dns   = "dns"
	yum   = "yum源"
	ali   = "ali"
	local = "local"
)

// todo
var (
	setHelpContent = "这是设置的帮助内容..."
	configFailed   = "配置失败..."
	configSuccess  = "配置成功..."
	successBanner  = "[success]"
	failedBanner   = "[failed]"

	Repo  string // 仓库地址
	Proxy string // 代理地址

	missingParameterErr = errors.New("缺少参数...")
)

func init() {
	setYumCmd.Flags().StringVarP(&Repo, "repo", "r", "", "Repository address of yum")
	setYumCmd.Flags().StringVarP(&Proxy, "proxy", "p", "", "Proxy address of yum")

	setCmd.AddCommand(setYumCmd)
	setCmd.AddCommand(setDNSCmd)
	setCmd.AddCommand(setHostnameCmd)
	setCmd.AddCommand(setTimeZoneCmd)

	rootCmd.AddCommand(setCmd)
	flag.Parse()

}

//// set 命令合法参数
var setValidArgs = []string{"dns", "yum", "hostname"}

// set命令
var setCmd = &cobra.Command{
	Use:   "set [OPTIONS] [flags]",
	Short: "set something through easyctl",
	Example: "\neasyctl set dns 114.114.114.114" +
		"\neasyctl set yum ali" +
		"\neasyctl set hostname weiliang.com",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(len(args))
	},
	ValidArgs: setValidArgs,
	Args:      cobra.ExactValidArgs(1),
}

// 配置yum子命令
var setYumCmd = &cobra.Command{
	Use:   "yum [flags]",
	Short: "easyctl set yum [flags]",
	Example: "\neasyctl set yum --repo=ali" +
		"\neasyctl set yum --repo=local" +
		"\neasyctl set yum --proxy=http://xxx:xxx@xxx.xxx.xxx.xxx:xxx",
	Run: func(cmd *cobra.Command, args []string) {
		setYUM()
	},
}

// 配置dns子命令
var setDNSCmd = &cobra.Command{
	Use:     "dns [value]",
	Short:   "easyctl set dns [value]",
	Example: "\neasyctl set dns 8.8.8.8",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		setDNS(args[0])
	},
}

// 配置hostname子命令
var setHostnameCmd = &cobra.Command{
	Use:     "hostname [value]",
	Short:   "easyctl set hostname [value]",
	Example: "\neasyctl set hostname nginx-server1",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		setHostname(args[0])
	},
}

// 配置时区子命令
var setTimeZoneCmd = &cobra.Command{
	Use:     "timezone",
	Short:   "easyctl set tz/timezone [value]",
	Example: "\neasyctl set tz/timezone",
	Aliases: []string{"tz"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		setTimeZone()
	},
}

// 配置dns
// todo 支持多dns地址
func setDNS(address string) {
	fmt.Printf("#### 配置dns地址：%s ####\n", address)
	err, _ := sys.SetDNS(address)
	if err != nil {
		fmt.Println(failedBanner + " " + dns + configFailed + ": " + err.Error())
	} else {
		fmt.Println(successBanner + dns + configSuccess)
	}
}

// 配置yum
// todo 待优化
func setYUM() {
	if Repo != "" && Repo == ali {
		sys.SetAliYUM()
	}
	if Repo != "" && Repo == local {
		sys.SetLocalYUM()
	}
}

// 配置hostname
func setHostname(name string) {
	sys.SetHostname(name)
}

// 配置时区
func setTimeZone() {
	sys.SetTimeZone()
}
