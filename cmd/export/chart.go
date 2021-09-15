package export

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/export"
)

// 导出chart
var chartCmd = &cobra.Command{
	Use:   "chart [flags]",
	Short: "export chart from harbor through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		export.Chart(serverListFile)
	},
}
