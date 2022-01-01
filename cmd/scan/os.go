package scan

import (
	//
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/scan"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

//go:embed asset/os.yaml
var osConfig []byte

var osCmd = &cobra.Command{
	Use:   "os [flags]",
	Short: "os命令指令集",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(
			command.Item{Cmd: cmd, Fnc: scan.OS, DefaultConfig: osConfig, ConfigFilePath: configFile}); runErr.Err != nil {
			command.DefaultLogger.Errorln(runErr.Msg)
			panic(runErr.Err)
		}
	},
}
