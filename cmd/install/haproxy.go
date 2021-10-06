package install

import (
	"github.com/spf13/cobra"
)

// install keepalive
var haproxyCmd = &cobra.Command{
	Use:   "haproxy [flags]",
	Short: "install haproxy through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
	},
	Args: cobra.NoArgs,
}
