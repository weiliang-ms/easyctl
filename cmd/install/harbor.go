package install

import (
	"github.com/spf13/cobra"
)

// install docker-compose
var harborCmd = &cobra.Command{
	Use:   "harbor [flags]",
	Short: "install harbor through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
	},
	Args: cobra.NoArgs,
}
