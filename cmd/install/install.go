package install

import (
	"github.com/spf13/cobra"
)

var configFile string

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")

	//RootCmd.AddCommand(keepaliveCmd)
	//RootCmd.AddCommand(haproxyCmd)
	RootCmd.AddCommand(dockerCmd)
	//RootCmd.AddCommand(dockerComposeCmd)
	//RootCmd.AddCommand(harborCmd)
	RootCmd.AddCommand(redisClusterCmd)
}

// RootCmd 安装根指令
var RootCmd = &cobra.Command{
	Use:   "install [OPTIONS] [flags]",
	Short: "安装指令集",
	Args:  cobra.ExactValidArgs(1),
}
