package export

import (
	"github.com/spf13/cobra"
)

var (
	packageName string
)

// export命令
var RootCmd = &cobra.Command{
	Use:     "export [OPTIONS] [flags]",
	Short:   "export something through easyctl",
	Example: "\nexport yum-repo --package-name=gcc",
	Run: func(cmd *cobra.Command, args []string) {
	},
	ValidArgs: []string{"yum-repo"},
	Args:      cobra.ExactValidArgs(1),
}

func init() {
	RootCmd.AddCommand(imageCmd)
}
