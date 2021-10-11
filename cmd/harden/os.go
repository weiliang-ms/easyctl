package harden

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/harden"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

// 安全加固命令
var osHardenCmd = &cobra.Command{
	Use: "os [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(command.Item{
			Cmd:            cmd,
			Fnc:            harden.OS,
			ConfigFilePath: configFile,
		}); runErr.Err != nil {
			panic(runErr.Err)
		}
	},
}
