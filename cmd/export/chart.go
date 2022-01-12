package export

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/export/chart"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"log"
)

//go:embed asset/chart-repo.yaml
var chartConfig []byte

// 导出chart
var chartCmd = &cobra.Command{
	Use:   "chart [flags]",
	Short: "导出charts指令",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(command.Item{
			Cmd:            cmd,
			DefaultConfig:  chartConfig,
			Fnc:            chart.Run,
			ConfigFilePath: configFile,
		}); runErr.Err != nil {
			log.Println(runErr.Msg)
			panic(runErr.Err)
		}
	},
}
