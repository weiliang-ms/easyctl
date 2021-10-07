package track

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/track"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

//go:embed asset/tail_log.yaml
var tailLogConfig []byte

// RootCmd close命令
var tailLogCmd = &cobra.Command{
	Use:   "tail-log [flags]",
	Short: "追踪日志命令",
	Run: func(cmd *cobra.Command, args []string) {
		if err := command.SetExecutorDefault(command.Item{
			Cmd:            cmd,
			Fnc:            track.TaiLog,
			DefaultConfig:  tailLogConfig,
			ConfigFilePath: configFile,
		}); err != nil {
			panic(err)
		}
	},
}
