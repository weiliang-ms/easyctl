package set

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/set"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

//go:embed asset/newpassword_config.yaml
var newPasswordConfig []byte

// 修改主机密码
var newPasswordCmd = &cobra.Command{
	Use:   "new-password [flags]",
	Short: "修改主机root口令",
	Args:  cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(command.Item{
			Cmd:            cmd,
			Fnc:            set.NewPassword,
			DefaultConfig:  newPasswordConfig,
			ConfigFilePath: configFile,
		}); runErr.Err != nil {
			panic(runErr.Err)
		}
	},
}
