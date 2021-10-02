package scan

import (
	"github.com/spf13/cobra"
)

// 配置hostname子命令
var scanOSCmd = &cobra.Command{
	Use:   "os [flags]",
	Short: "easyctl scan os",
	Run: func(cmd *cobra.Command, args []string) {
		//scan.OSSecurity()
	},
}
