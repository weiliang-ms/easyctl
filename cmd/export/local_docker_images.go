package export

import (
	"github.com/spf13/cobra"
)

// install docker-compose
var localImagesCmd = &cobra.Command{
	Use:   "local-image-list",
	Short: "导出本地镜像列表",
	Run: func(cmd *cobra.Command, args []string) {
		//localImageList()
	},
	Args: cobra.NoArgs,
}
