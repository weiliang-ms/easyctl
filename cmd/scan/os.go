package scan

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/scan"
)

// 配置hostname子命令
var scanOSCmd = &cobra.Command{
	Use:   "os [flags]",
	Short: "easyctl scan os",
	Run: func(cmd *cobra.Command, args []string) {
		scan.OSSecurity()
	},
}
