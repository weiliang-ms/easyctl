package set

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/set"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"log"
)

// 配置时区子命令
var timeZoneCmd = &cobra.Command{
	Use:     "timezone",
	Short:   "设置为上海时区",
	Example: "\neasyctl set tz/timezone",
	Aliases: []string{"tz"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := command.SetExecutorDefault(command.Item{
			Cmd:            cmd,
			Fnc:            set.Timezone,
			ConfigFilePath: configFile,
		}); err.Err != nil {
			log.Println(err.Msg)
			panic(err.Err)
		}
	},
}
