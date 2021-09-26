package set

import (
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/set"
)

//go:embed asset/dns_config.yaml
var dnsConfig []byte

// 配置dns子命令
var dnsCmd = &cobra.Command{
	Use:   "dns [flags]",
	Short: "配置主机dns",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := Set(Entity{Cmd: cmd, Fnc: set.Dns, DefaultConfig: dnsConfig}); runErr != nil {
			panic(runErr)
		}
	},
}
