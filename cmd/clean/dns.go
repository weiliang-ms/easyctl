package clean

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/clean"
)

// clean命令
var dnsCmd = &cobra.Command{
	Use:     "dns [flags]",
	Short:   "清理dns",
	Example: "\neasyctl clean dns",
	Run: func(cmd *cobra.Command, args []string) {
		cleanDns()
	},
}

func cleanDns() {
	clean := clean.Actuator{
		ServerListFile: serverListFile,
	}
	clean.Dns()
}
