package upgrade

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(upgradeKernelCmd)
}

// upgrade 命令
var Cmd = &cobra.Command{
	Use:     "upgrade [OPTIONS] [flags]",
	Short:   "upgrade soft through easyctl",
	Example: "\neasyctl upgrade kernel --kernel-version=lt\n",
	Args:    cobra.MinimumNArgs(1),
}
