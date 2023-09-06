package exec

import (
	// embed
	_ "embed"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"os"
)

var (
	configFile string
)

//go:embed asset/exec.yaml
var execConfig []byte

// RootCmd close命令
var RootCmd = &cobra.Command{
	Use:     "exec [flags]",
	Short:   "执行命令指令集",
	Example: "\neasyctl exec 'date' -c config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			shell := args[0]
			command.DefaultLogger.Infof("执行: %s", shell)
			b, err := os.ReadFile(configFile)
			if err != nil {
				logrus.Infof("未找到配置文件，为您生成配置文件样例, 请修改文件内容后携带 -c 参数重新执行 -> %s", constant.ConfigFile)
				_ = os.WriteFile(constant.ConfigFile, execConfig, 0666)
				os.Exit(-1)
			}
			runner.RemoteRun(runner.RemoteRunItem{
				ManifestContent:     b,
				Logger:              command.DefaultLogger,
				Cmd:                 shell,
				RecordErrServerList: true,
			})
		}
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(shellCmd)
	RootCmd.AddCommand(suShellCmd)
	RootCmd.AddCommand(scpCmd)
	RootCmd.AddCommand(pingCmd)
}
