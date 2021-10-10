package install

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/install"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

//go:embed asset/redis_cluster.yaml
var redisClusterConfig []byte

// redisClusterCmd 安装redis集群指令
var redisClusterCmd = &cobra.Command{
	Use:   "redis-cluster [flags]",
	Short: "安装redis集群",
	Args:  cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(
			command.Item{
				Cmd:            cmd,
				Fnc:            install.RedisCluster,
				DefaultConfig:  redisClusterConfig,
				ConfigFilePath: configFile,
			}); runErr.Err != nil {
			panic(runErr.Err)
		}
	},
}
