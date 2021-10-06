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
		if err := command.SetExecutorDefault(command.ExecutorEntity{
			Cmd:           cmd,
			Fnc:           set.NewPassword,
			DefaultConfig: newPasswordConfig,
		}, configFile); err != nil {
			panic(err)
		}
	},
}
