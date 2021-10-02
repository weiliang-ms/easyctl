package export

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/export"
)

//go:embed asset/chart-repo.yaml
var chartConfig []byte

// 导出chart
var chartCmd = &cobra.Command{
	Use:   "chart [flags]",
	Short: "导出charts指令",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := Export(Entity{Cmd: cmd, Fnc: export.Chart, DefaultConfig: chartConfig}); runErr != nil {
			panic(runErr)
		}
	},
}
