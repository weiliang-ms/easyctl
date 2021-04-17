package install

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/api/install"
	"github.com/weiliang-ms/easyctl/pkg/runner"
)

func init() {
	haproxyCmd.Flags().BoolVarP(&offline, "offline", "", false, "是否离线安装")
	haproxyCmd.Flags().StringVarP(&offlineFilePath, "offline-file", "", "", "离线文件")
	haproxyCmd.Flags().StringVarP(&serverListFile, "server-list", "", "", "服务器批量连接信息")
}

// install keepalive
var haproxyCmd = &cobra.Command{
	Use:   "haproxy [flags]",
	Short: "install haproxy through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		haproxy()
	},
	Args: cobra.NoArgs,
}

func haproxy() {
	i := runner.Installer{
		ServerListPath:  serverListFile,
		Offline:         offline,
		OfflineFilePath: offlineFilePath,
	}
	install.Haproxy(i)
}
