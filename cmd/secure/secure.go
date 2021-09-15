package secure

import (
	"github.com/spf13/cobra"
)

// 安全加固命令
var RootCmd = &cobra.Command{
	Use:   "secure [OPTIONS] [flags]",
	Short: "secure something through easyctl",
	Args:  cobra.ExactValidArgs(1),
}

func init() {
	//RootCmd.AddCommand(osSecureCmd)
}
