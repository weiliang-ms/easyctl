package clean

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/clean"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

//go:embed asset/dns_config.yaml
var dnsDefaultConfig []byte

// clean命令
var cleanDnsCmd = &cobra.Command{
	Use:     "dns [flags]",
	Short:   "清理dns列表",
	Example: "\neasyctl clean dns",
	Run: func(cmd *cobra.Command, args []string) {
		if err := command.SetExecutorDefault(command.ExecutorEntity{
			Cmd:           cmd,
			Fnc:           clean.Dns,
			DefaultConfig: dnsDefaultConfig,
		}, configFile); err != nil {
			panic(err)
		}
	},
}
