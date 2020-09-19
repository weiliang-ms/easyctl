package cmd

import (
	"easyctl/sys"
	"github.com/spf13/cobra"
)

var Nologin bool

func init() {

	addUserCmd.Flags().BoolVarP(&Nologin, "no-login", "n", false, "User type: no login")
	addUserCmd.Flags().Parsed()

	addCmd.AddCommand(addUserCmd)
	rootCmd.AddCommand(addCmd)

}

// add命令
var addCmd = &cobra.Command{
	Use:     "add [OPTIONS] [flags]",
	Short:   "add something through easyctl",
	Example: "\neasyctl add user user1 password",
	Run: func(cmd *cobra.Command, args []string) {
	},
	ValidArgs: []string{"user"},
	Args:      cobra.ExactValidArgs(1),
}

// addUser命令
var addUserCmd = &cobra.Command{
	Use:   "user [username] [password] [flags]",
	Short: "add linux user through easyctl, password default value: user123",
	Example: "\neasyctl add user user1 password" +
		"\neasyctl add user user1 password --no-login",
	Run: func(cmd *cobra.Command, args []string) {
		addUser(args)
	},
	Args: cobra.MinimumNArgs(1),
}

func addUser(args []string) {
	password := ""
	username := args[0]
	if len(args) > 1 {
		password = args[1]
	}
	sys.AddUser(username, password, !Nologin)
}
