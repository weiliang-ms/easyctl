package cmd

import (
	"github.com/spf13/cobra"
)

func init() {

}

// install docker
var installDockerCmd = &cobra.Command{
	Use:   "docker [flags]",
	Short: "install docker through easyctl...",
	// Example: "\neasyctl download harbor --url=https://github.com/goharbor/harbor/releases/download/v2.1.4/harbor-offline-installer-v2.1.4.tgz",
	Run: func(cmd *cobra.Command, args []string) {
		download(Url, "harbor")
	},
	Args: cobra.NoArgs,
}

func setLocalYum() {

}
