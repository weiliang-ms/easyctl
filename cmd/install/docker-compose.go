package install

import (
	"github.com/spf13/cobra"
)

// install docker-compose
var dockerComposeCmd = &cobra.Command{
	Use:   "docker-compose [flags]",
	Short: "install docker-compose through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
	},
	Args: cobra.NoArgs,
}
