package scan

import (
	"github.com/spf13/cobra"
)

// scan命令
var RootCmd = &cobra.Command{
	Use:   "scan [OPTIONS]",
	Short: "scan something through easyctl",
}

func init() {
	RootCmd.AddCommand(scanOSCmd)
}
