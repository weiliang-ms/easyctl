package cmd

import "github.com/spf13/cobra"

func init() {
	upgradeCmd.AddCommand(upgradeKernelCmd)
}

// upgrade 命令
var upgradeCmd = &cobra.Command{
	Use:     "upgrade [OPTIONS] [flags]",
	Short:   "upgrade soft through easyctl",
	Example: "\neasyctl upgrade kernel --kernel-version=lt\n",
	Args:    cobra.MinimumNArgs(1),
}
