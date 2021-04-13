package close

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctll/pkg/runner"
)

var (
	forever bool
	remote bool
	serverListFile string
)

func init() {
	RootCmd.AddCommand(closeFirewallCmd)
	RootCmd.AddCommand(closeSeLinuxCmd)
}

// close命令
var RootCmd = &cobra.Command{
	Use:   "close [OPTIONS] [flags]",
	Short: "close some service through easyctl",
	Example: "\neasyctl close firewalld",
	Run: func(cmd *cobra.Command, args []string) {
	},
	ValidArgs: []string{"firewall"},
	Args:      cobra.ExactValidArgs(1),
}


func close(cmd string,list []runner.Server)  {
	if len(list) == 0 {
		runner.Shell(cmd)
		return
	}

	for _, v := range list{
		v.RemoteShell(cmd)
	}
}
