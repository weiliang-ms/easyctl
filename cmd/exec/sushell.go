package exec

import (
	// _
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/exec"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

//go:embed asset/su-executor.yaml
var suShellConfig []byte

var suShellCmd = &cobra.Command{
	Use:   "su-shell [flags]",
	Short: "切换用户执行命令指令集",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(
			command.Item{Cmd: cmd, Fnc: exec.SURun, DefaultConfig: suShellConfig, ConfigFilePath: configFile}); runErr.Err != nil {
			command.DefaultLogger.Errorln(runErr.Msg)
			panic(runErr.Err)
		}
	},
}
