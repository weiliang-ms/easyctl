package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Url string
)

func init() {

	downloadCmd.AddCommand(downloadHarborCmd)
	RootCmd.AddCommand(downloadCmd)

}

// add命令
var downloadCmd = &cobra.Command{
	Use:   "download [OPTIONS] [flags]",
	Short: "download soft through easyctl",
	Run: func(cmd *cobra.Command, args []string) {
	},
	ValidArgs: []string{"harbor"},
	Args:      cobra.ExactValidArgs(1),
}
