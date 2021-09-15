package add

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/add"
)

func init() {
	addUserCmd.Flags().BoolVarP(&Nologin, "no-login", "n", false, "非登录用户")
	addUserCmd.Flags().StringVarP(&username, "username", "u", "", "用户名")
	addUserCmd.Flags().StringVarP(&password, "password", "p", "", "用户密码")
	_ = addUserCmd.MarkFlagRequired("username")
}

// addUser命令
var addUserCmd = &cobra.Command{
	Use:   "user [flags]",
	Short: "add linux user through easyctl, password default value: user123",
	Example: "\neasyctl add user -u user1 -p password" +
		"\neasyctl add user -u user1 --no-login",
	Run: func(cmd *cobra.Command, args []string) {
		addUser()
	},
	Args: cobra.NoArgs,
}

func addUser() {
	ac := add.Actuator{
		ServerListFile: serverListFile,
		UserName:       username,
		Password:       password,
		NoLogin:        Nologin,
	}
	ac.AddUser()
}
