package clean

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/clean"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"log"
)

// 清理redis命令
var cleanRedisCmd = &cobra.Command{
	Use:   "redis [flags]",
	Short: "清理redis文件及服务",
	Run: func(cmd *cobra.Command, args []string) {
		if err := command.SetExecutorDefault(
			command.Item{
				Cmd:            cmd,
				Fnc:            clean.Redis,
				ConfigFilePath: configFile,
			}); err.Err != nil {
			log.Println(err.Msg)
			panic(err.Err)
		}
	},
}
