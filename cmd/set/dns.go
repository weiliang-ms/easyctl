package set

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/set"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"log"
)

//go:embed asset/dns_config.yaml
var dnsConfig []byte

// 配置dns子命令
var dnsCmd = &cobra.Command{
	Use:   "dns [flags]",
	Short: "配置主机dns",
	Run: func(cmd *cobra.Command, args []string) {
		if err := command.SetExecutorDefault(command.Item{
			Cmd:            cmd,
			Fnc:            set.Dns,
			DefaultConfig:  dnsConfig,
			ConfigFilePath: configFile,
		}); err.Err != nil {
			log.Println(err.Msg)
			panic(err.Err)
		}
	},
}
