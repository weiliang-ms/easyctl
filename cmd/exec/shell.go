package exec

import (
	// _
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/exec"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

//go:embed asset/executor.yaml
var shellConfig []byte

var shellCmd = &cobra.Command{
	Use:     "shell [flags]",
	Short:   "执行命令指令集",
	Example: "\neasyctl exec -c config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(
			command.Item{Cmd: cmd, Fnc: exec.Run, DefaultConfig: shellConfig, ConfigFilePath: configFile}); runErr.Err != nil {
			command.DefaultLogger.Errorln(runErr.Msg)
			panic(runErr.Err)
		}
	},
}
