package harden

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/harden"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"log"
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
			log.Println(runErr.Msg)
			panic(runErr.Err)
		}
	},
}
