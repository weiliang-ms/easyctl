package install

import (
	"github.com/spf13/cobra"
)

// install docker
var dockerCmd = &cobra.Command{
	Use:   "docker-ce [flags]",
	Short: "install docker through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
