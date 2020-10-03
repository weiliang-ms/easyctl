package cmd

import (
	"easyctl/sys"
	"github.com/spf13/cobra"
)

var (
	Nologin  bool
	username string
	password string
)

func init() {

	addUserCmd.Flags().BoolVarP(&Nologin, "no-login", "n", false, "User type: no login")
	addUserCmd.Flags().StringVarP(&username, "username", "u", "", "user name")
	addUserCmd.Flags().StringVarP(&password, "password", "p", "", "user password")
	addUserCmd.MarkFlagRequired("username")

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
	if password == "" {
		password = "user123"
	}
	sys.AddUser(username, password, !Nologin)
}
