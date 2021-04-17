package install

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/api/install"
	"github.com/weiliang-ms/easyctl/pkg/runner"
)

func init() {
	dockerCmd.Flags().BoolVarP(&offline, "offline", "", false, "是否离线安装")
	dockerCmd.Flags().StringVarP(&offlineFilePath, "offline-file", "", "", "离线文件")
	dockerCmd.Flags().StringVarP(&serverListFile, "server-list", "", "", "服务器批量连接信息")
}

// install docker
var dockerCmd = &cobra.Command{
	Use:   "docker-ce [flags]",
	Short: "install docker through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		docker()
	},
}

// 单机本地离线
func docker() {
	i := runner.Installer{
		ServerListPath:  serverListFile,
		Offline:         offline,
		OfflineFilePath: offlineFilePath,
	}
	install.Docker(i)
}
