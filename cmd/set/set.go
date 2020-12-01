package set

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/sys"
)

// todo
var (
	CheckIDRsaFileCmd    = fmt.Sprintf("[ -e $HOME/.ssh/%s ]", idRsa)
	CheckIDRsaPubFileCmd = fmt.Sprintf("[ -e $HOME/.ssh/%s ]", idRsaPub)

	HostsFilePath string // 主机互信hosts文件路径

	missingParameterErr = errors.New("缺少参数...")
)

//// set 命令合法参数
var setValidArgs = []string{"dns", "yum", "hostname", "timezone", "pubkey-authentication"}

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
	//err, _ := sys.SetDNS(address)
	//if err != nil {
	//	fmt.Println(failedBanner + " " + dns + configFailed + ": " + err.Error())
	//} else {
	//	fmt.Println(successBanner + dns + configSuccess)
	//}
}

// 配置hostname
func setHostname(name string) {
	sys.SetHostname(name)
}

// 配置时区
func setTimeZone() {
	sys.SetTimeZone()
}
