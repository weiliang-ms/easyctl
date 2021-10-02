package add

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/add"
)

//go:embed asset/user.yaml
var userConfig []byte

// addUser命令
var addUserCmd = &cobra.Command{
	Use:   "user [flags]",
	Short: "创建用户指令",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := Add(Entity{Cmd: cmd, Fnc: add.User, DefaultConfig: userConfig}); runErr != nil {
			panic(runErr)
		}
	},
	Args: cobra.NoArgs,
}
