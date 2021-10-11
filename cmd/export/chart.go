package export

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/export"
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
		options := make(map[string]interface{})
		options[export.GetChartListFunc] = export.GetChartList
		options[export.GetChartsByteFunc] = export.GetChartsByte

		if runErr := command.SetExecutorDefault(command.Item{
			Cmd:            cmd,
			DefaultConfig:  chartConfig,
			Fnc:            export.Chart,
			ConfigFilePath: configFile,
			OptionFunc:     options,
		}); runErr.Err != nil {
			log.Println(runErr.Msg)
			panic(runErr.Err)
		}
	},
}
