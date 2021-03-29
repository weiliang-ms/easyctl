package cmd

import "github.com/spf13/cobra"

func init() {
	RootCmd.AddCommand(installCmd)
}

// install命令
var installCmd = &cobra.Command{
	Use:     "install [OPTIONS] [flags]",
	Short:   "install soft through easyctl",
	Example: "\neasyctl install docker\n",
	Args:    cobra.MinimumNArgs(1),
}
