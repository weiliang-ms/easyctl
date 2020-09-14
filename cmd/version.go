package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// 输出easycfg版本
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of easycfg",
	Long:  `All software has versions. This is easycfg's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("easycfg  v0.1 -- alpha")
	},
}
