package set

import (
	"github.com/spf13/cobra"
)

// 配置dns子命令
// todo
var dnsCmd = &cobra.Command{
	Use:   "dns [flags]",
	Short: "easyctl set dns",
	Run: func(cmd *cobra.Command, args []string) {
		//setDNS()
	},
}

// 配置dns
//func setDNS() {
//	ac := &set.Actuator{
//		ServerListFile: serverListFile,
//		Value:          value,
//	}
//
//	ac.DNS()
//}
