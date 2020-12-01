package set

import (
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

const (
	dns        = "dns"
	ali        = "ali"
	local      = "local"
	idRsa      = "id_rsa"
	ak         = "authorized_keys"
	idRsaPub   = "id_rsa.pub"
	serverList = "server-list"
	isoPath    = "iso-path"
)

var (
	yumServerListFile   string
	trustServerListFile string // 主机互信主机列表
	imageFilePath       string
	yumRepo             string // 仓库地址
	yumProxy            string // 代理地址
)

// set命令
var RootCmd = &cobra.Command{
	Use:   "set [OPTIONS] [flags]",
	Short: "set something through easyctl",
	Example: "\neasyctl set dns 114.114.114.114" +
		"\neasyctl set yum ali" +
		"\neasyctl set hostname weiliang.com",
	Args: cobra.ExactValidArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cmd.ValidateArgs(setValidArgs)
	},
	ValidArgs: setValidArgs,
}

func init() {
	setYumCmd.Flags().StringVarP(&yumRepo, "repo", "r", "", "Repository address of yum")
	setYumCmd.Flags().StringVarP(&yumProxy, "proxy", "p", "", "Proxy address of yum")
	setYumCmd.Flags().StringVarP(&yumServerListFile, serverList, "", "", "配置yum主机列表")
	setYumCmd.Flags().StringVarP(&imageFilePath, isoPath, "", "", "本机系统版本镜像路径")

	setPubKeyAuthenticationCmd.Flags().StringVarP(&trustServerListFile, serverList, "", "", "配置主机互信的主机列表")
	setPubKeyAuthenticationCmd.MarkFlagRequired(serverList)

	RootCmd.AddCommand(setYumCmd)
	RootCmd.AddCommand(setDNSCmd)
	RootCmd.AddCommand(setHostnameCmd)
	RootCmd.AddCommand(setTimeZoneCmd)
	RootCmd.AddCommand(setPubKeyAuthenticationCmd)
	flag.Parse()

}

// 配置yum子命令
var setYumCmd = &cobra.Command{
	Use:   "yum [flags]",
	Short: "easyctl set yum [flags]",
	Example: "\neasyctl set yum --repo=ali" +
		"\neasyctl set yum --repo=local" +
		"\neasyctl set yum --proxy=http://username:password@192.168.111.222:8080",
	Args: cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		setYum(cmd)
	},
}

// 主机互信
var setPubKeyAuthenticationCmd = &cobra.Command{
	Use:     "pubkey-authentication [flags]",
	Short:   "easyctl set pubkey-authentication --hosts-file-path=./hosts.txt",
	Example: "\neasyctl set pubkey-authentication --hosts-file-path=./hosts.txt",
	Aliases: []string{"pka"},
	Args:    cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var trust trust
		trust.setTrust()
	},
}

func needFlag(cmd *cobra.Command) {
	if cmd.Flags().NFlag() == 0 {
		fmt.Printf("Flags:\n%s", cmd.Flags().FlagUsages())
	}
}

func setYum(cmd *cobra.Command) {
	needFlag(cmd)
	if yumRepo == local && imageFilePath == "" {
		log.Fatal("配置本地yum源，必须通过--iso-path指定iso镜像路径")
	}
	var yum yum
	yum.setYum()
}
