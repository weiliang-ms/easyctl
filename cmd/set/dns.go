package set

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/set"
)

//

func init() {
	dnsCmd.Flags().StringVarP(&value, "value", "v", "", "dns地址，多个地址以 ','隔离 ")
	dnsCmd.MarkFlagRequired("value")
}

// 配置dns子命令
var dnsCmd = &cobra.Command{
	Use:   "dns [flags]",
	Short: "easyctl set dns --value",
	Example: "\neasyctl set dns --value=8.8.8.8\n" +
		"easyctl set dns --value=8.8.8.8,114.114.114.114",
	Run: func(cmd *cobra.Command, args []string) {
		setDNS()
	},
}

// 配置dns
func setDNS() {
	ac := &set.Actuator{
		ServerListFile: serverListFile,
		Value:          value,
	}

	ac.DNS()
}
