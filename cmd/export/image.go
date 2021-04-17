package export

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/api/export"
)

var serverListFile string

func init() {
	imageCmd.Flags().StringVarP(&serverListFile, "server-list", "", "", "服务器批量连接信息")
	imageCmd.MarkFlagRequired("domain")
}

// install docker-compose
var imageCmd = &cobra.Command{
	Use:   "image [flags]",
	Short: "export harbor images through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		exportImage()
	},
	Args: cobra.NoArgs,
}

// 单机本地离线
func exportImage() {
	export.HarborImage(serverListFile)
}
