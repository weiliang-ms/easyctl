package clean

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/clean"
)

//go:embed asset/dns_config.yaml
var dnsDefaultConfig []byte

// clean命令
var cleanDnsCmd = &cobra.Command{
	Use:     "dns [flags]",
	Short:   "清理dns列表",
	Example: "\neasyctl clean dns",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := Clean(Entity{Cmd: cmd, Fnc: clean.Dns, DefaultConfig: dnsDefaultConfig}); runErr != nil {
			panic(runErr)
		}
	},
}
