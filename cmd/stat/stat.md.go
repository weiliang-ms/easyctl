package stat

import (
	"github.com/spf13/cobra"
)

// RootCmd stat命令
var RootCmd = &cobra.Command{
	Use:   "stat [OPTIONS] [flags]",
	Short: "get system settings through easyctl",
}
