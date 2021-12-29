package exec

import (
	//
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/exec"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

//go:embed asset/ping.yaml
var pingConfig []byte

var pingCmd = &cobra.Command{
	Use:     "ping [flags]",
	Short:   "ping命令指令集",
	Example: "\neasyctl exec -c config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(
			command.Item{Cmd: cmd, Fnc: exec.Ping, DefaultConfig: pingConfig, ConfigFilePath: configFile}); runErr.Err != nil {
			command.DefaultLogger.Errorln(runErr.Msg)
			panic(runErr.Err)
		}
	},
}
