package harden

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/harden"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"log"
)

func init() {
	osHardenCmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "Skip confirm step.")
	osHardenCmd.Flags().BoolVarP(&localRun, "local", "", false, "harden current server.")
}

// 安全加固命令
var osHardenCmd = &cobra.Command{
	Use: "os [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(command.Item{
			Cmd:            cmd,
			Fnc:            harden.OS,
			ConfigFilePath: configFile,
			SkipConfirm:    skipConfirm,
			DefaultConfig:  config,
			LocalRun:       localRun,
		}); runErr.Err != nil {
			log.Println(runErr.Msg)
			panic(runErr.Err)
		}
	},
}
