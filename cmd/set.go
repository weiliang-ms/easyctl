package cmd

import (
	"easyctl/sys"
	"easyctl/util"
	"errors"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
)

const (
	dns                         = "dns"
	ali                         = "ali"
	local                       = "local"
	idRsa                       = "id_rsa"
	idRsaPub                    = "id_rsa.pub"
	authorizedKeys              = "authorized_keys"
	generateIDRsaAndIDRsaPubCmd = "ssh-keygen -t rsa -N '' -f ~/.ssh/id_rsa -q"
)

// todo
var (
	setHelpContent = "这是设置的帮助内容..."
	configFailed   = "配置失败..."
	configSuccess  = "配置成功..."
	successBanner  = "[success]"
	failedBanner   = "[failed]"

	CheckIDRsaFileCmd    = fmt.Sprintf("[ -e $HOME/.ssh/%s ]", idRsa)
	CheckIDRsaPubFileCmd = fmt.Sprintf("[ -e $HOME/.ssh/%s ]", idRsaPub)

	Repo          string // 仓库地址
	Proxy         string // 代理地址
	HostsFilePath string // 主机互信hosts文件路径

	missingParameterErr = errors.New("缺少参数...")
)

func init() {
	setYumCmd.Flags().StringVarP(&Repo, "repo", "r", "", "Repository address of yum")
	setYumCmd.Flags().StringVarP(&Proxy, "proxy", "p", "", "Proxy address of yum")
	setPubKeyAuthenticationCmd.Flags().StringVarP(&HostsFilePath, "hosts-file-path", "", "", "配置主机互信的主机列表文件")
	setPubKeyAuthenticationCmd.MarkFlagRequired("hosts-file-path")

	setCmd.AddCommand(setYumCmd)
	setCmd.AddCommand(setDNSCmd)
	setCmd.AddCommand(setHostnameCmd)
	setCmd.AddCommand(setTimeZoneCmd)
	setCmd.AddCommand(setPubKeyAuthenticationCmd)

	rootCmd.AddCommand(setCmd)
	flag.Parse()

}

//// set 命令合法参数
var setValidArgs = []string{"dns", "yum", "hostname", "timezone", "pubkey-authentication"}

// set命令
var setCmd = &cobra.Command{
	Use:   "set [OPTIONS] [flags]",
	Short: "set something through easyctl",
	Example: "\neasyctl set dns 114.114.114.114" +
		"\neasyctl set yum ali" +
		"\neasyctl set hostname weiliang.com",
	RunE: func(cmd *cobra.Command, args []string) error {
		return parseCommand(cmd, args, setValidArgs)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return setValidArgs, cobra.ShellCompDirectiveNoFileComp
	},
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

// 主机互信
var setPubKeyAuthenticationCmd = &cobra.Command{
	Use:     "pubkey-authentication [flags]",
	Short:   "easyctl set pubkey-authentication --hosts-file-path=./hosts.txt",
	Example: "\neasyctl set pubkey-authentication --hosts-file-path=./hosts.txt",
	Aliases: []string{"pka"},
	Run: func(cmd *cobra.Command, args []string) {
		setServerTrust()
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

// 设置主机互信
// todo 待优化代码
func setServerTrust() {
	// 1.读取ssh对象数组
	sshObjects := util.ReadSSHInfoFromFile(HostsFilePath)
	if len(sshObjects) < 1 {
		panic("ssh信息禁止为空...")
	}

	// 2.读取列表第一个主机.ssh信息
	object := sshObjects[0]
	_, rsa := object.ExecuteOriginCmd(CheckIDRsaFileCmd, 0)
	_, rsaPub := object.ExecuteOriginCmd(fmt.Sprintf("[ -e $HOME/.ssh/%s ]", idRsaPub), 0)
	if rsa != 0 || rsaPub != 0 {
		object.ExecuteOriginCmd(generateIDRsaAndIDRsaPubCmd, 0)
	}

	// 3.获取文件内容
	rsaContent, _ := object.ExecuteOriginCmd(fmt.Sprintf("cat $HOME/.ssh/%s", idRsa), 0)
	rsaPubContent, _ := object.ExecuteOriginCmd(fmt.Sprintf("cat $HOME/.ssh/%s", idRsaPub), 0)
	fmt.Printf("rsa:\n%s \n\nrsapub:\n%s", rsaContent, rsaPubContent)

	// 4.自互信
	object.ExecuteOriginCmd(fmt.Sprintf("\\cp $HOME/.ssh/%s $HOME/.ssh/%s", idRsaPub, authorizedKeys), 0)
	fmt.Println("----后续逻辑...")
}
