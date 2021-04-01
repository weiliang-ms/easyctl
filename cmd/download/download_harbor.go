package download

import (
	"github.com/spf13/cobra"
)

func init() {
	downloadHarborCmd.Flags().StringVarP(&Url, "url", "",
		// todo: 读取配置文件
		"https://github.com/goharbor/harbor/releases/download/v2.1.4/harbor-offline-installer-v2.1.4.tgz",
		"harbor resource url")
}

// download harbor resources
var downloadHarborCmd = &cobra.Command{
	Use:     "harbor [flags]",
	Short:   "download linux soft resources through easyctl...",
	Example: "\neasyctl download harbor --url=https://github.com/goharbor/harbor/releases/download/v2.1.4/harbor-offline-installer-v2.1.4.tgz",
	Run: func(cmd *cobra.Command, args []string) {
		download(Url, "harbor")
	},
	Args: cobra.NoArgs,
}
