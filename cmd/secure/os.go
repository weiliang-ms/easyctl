package secure

import "github.com/spf13/cobra"

// 安全加固命令
var osSecureCmd = &cobra.Command{
	Use:   "os [flags]",
	Short: "secure os through easyctl",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}
