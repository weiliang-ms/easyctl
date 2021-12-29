package exec

import (
	//
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/exec"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

//go:embed asset/scp.yaml
var scpConfig []byte

var scpCmd = &cobra.Command{
	Use:     "scp [flags]",
	Short:   "拷贝命令指令集",
	Example: "\neasyctl exec -c config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(
			command.Item{Cmd: cmd, Fnc: exec.Run, DefaultConfig: scpConfig, ConfigFilePath: configFile}); runErr.Err != nil {
			command.DefaultLogger.Errorln(runErr.Msg)
			panic(runErr.Err)
		}
	},
}
