package export

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/export"
)

// install docker-compose
var localImagesCmd = &cobra.Command{
	Use:   "local-image-list",
	Short: "导出本地镜像列表",
	Run: func(cmd *cobra.Command, args []string) {
		localImageList()
	},
	Args: cobra.NoArgs,
}

// 单机本地离线
func localImageList() {
	export.LocalImageList()
}
