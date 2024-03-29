package add

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/add"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"log"
)

//go:embed asset/user.yaml
var userConfig []byte

// addUser命令
var addUserCmd = &cobra.Command{
	Use:   "user [flags]",
	Short: "创建用户指令",
	Run: func(cmd *cobra.Command, args []string) {
		if err := command.SetExecutorDefault(
			command.Item{
				Cmd:            cmd,
				Fnc:            add.User,
				DefaultConfig:  userConfig,
				ConfigFilePath: configFile,
			}); err.Err != nil {
			log.Println(err.Msg)
			panic(err.Err)
		}
	},
	Args: cobra.NoArgs,
}
