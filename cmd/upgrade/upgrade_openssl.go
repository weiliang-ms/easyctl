package upgrade

import (
	"github.com/spf13/cobra"
)

func init() {
	upgradeOpensslCmd.Flags().StringVarP(&filePath, "file-path", "f", "", "文件路径或url")
	_ = upgradeOpensslCmd.MarkFlagRequired("file-path")
}

// install kernel
var upgradeOpensslCmd = &cobra.Command{
	Use:   "openssl [flags]",
	Short: "upgrade openssl through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
	},
	Args: cobra.NoArgs,
}
