package boot

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/boot"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"log"
)

//go:embed asset/config.yaml
var configByte []byte

var appWithCGroupsCmd = &cobra.Command{
	Use:   "app-with-cgroups [flags]",
	Short: "以配额限制的方式启动应用",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(
			command.Item{Cmd: cmd, Fnc: boot.AppWithCGroups, DefaultConfig: configByte, ConfigFilePath: configFile}); runErr.Err != nil {
			log.Println(runErr.Msg)
			panic(runErr.Err)
		}
	},
}
