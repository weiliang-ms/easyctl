package install

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/api/install"
	"github.com/weiliang-ms/easyctl/pkg/runner"
)

var offlineImages string

func init() {
	harborCmd.Flags().BoolVarP(&offline, "offline", "", false, "是否离线安装")
	harborCmd.Flags().StringVarP(&offlineFilePath, "offline-file", "", "", "离线文件")
	harborCmd.Flags().StringVarP(&serverListFile, "server-list", "", "", "服务器批量连接信息")
	harborCmd.Flags().StringVarP(&offlineImages, "offline-images-path", "", "", "镜像文件")
	harborCmd.MarkFlagRequired("domain")
}

// install docker-compose
var harborCmd = &cobra.Command{
	Use:   "harbor [flags]",
	Short: "install harbor through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		harbor()
	},
	Args: cobra.NoArgs,
}

// 单机本地离线
func harbor() {
	i := runner.Installer{
		ServerListPath:  serverListFile,
		InitImagesPath:  offlineImages,
		Offline:         offline,
		OfflineFilePath: offlineFilePath,
	}
	install.Harbor(i)
}
