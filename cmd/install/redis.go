package install

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/install"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"log"
)

//go:embed asset/redis.yaml
var redisConfig []byte

// redisClusterCmd 安装redis指令
var redisCmd = &cobra.Command{
	Use:   "redis [flags]",
	Short: "安装redis",
	Args:  cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(
			command.Item{
				Cmd:            cmd,
				Fnc:            install.Redis,
				DefaultConfig:  redisClusterConfig,
				ConfigFilePath: configFile,
			}); runErr.Err != nil {
			log.Println(runErr.Msg)
			panic(runErr.Err)
		}
	},
}
