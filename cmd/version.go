package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// 输出easyctl版本
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of easyctl",
	Long:  `All software has versions. This is easyctl's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("easyctl  v0.1.0 -- alpha")
	},
}
