package set

import (
	"easyctl/asset"
	"easyctl/pkg/runner"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	setDNSCmd.Flags().StringVarP(&value, "value", "v", "", "dns 地址")
	setDNSCmd.Flags().BoolVarP(&multiNode, "multi-node", "", false, "是否配置多节点")
	setDNSCmd.Flags().StringVarP(&serverListFile, "server-list", "", "server.yaml", "服务器列表")
	setDNSCmd.MarkFlagRequired("value")
}

// 配置dns子命令
var setDNSCmd = &cobra.Command{
	Use:     "dns [flags]",
	Short:   "easyctl set dns --value",
	Example: "\neasyctl set dns --value=8.8.8.8",
	Run: func(cmd *cobra.Command, args []string) {
		setDNS()
	},
}

// 配置dns

func setDNS() {
	if !multiNode {
		setLocalDNS()
	} else {
		setMultiNodeDNS()
	}
}

func setLocalDNS() {
	local(fmt.Sprintf("配置dns,地址为%s...\n", value), dnsScript())
}

func setMultiNodeDNS() {
	list := runner.ParseServerList(serverListFile)
	multiShell(list, dnsScript())
}

func dnsScript() string {
	script, _ := asset.Asset("static/script/set/dns.sh")
	return fmt.Sprintf("address=%s %s", value, string(script))
}
