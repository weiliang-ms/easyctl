package install

import (
	"github.com/spf13/cobra"
)

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
	//i := runner.Installer{
	//	ServerListPath:  serverListFile,
	//	Offline:         offline,
	//	OfflineFilePath: offlineFilePath,
	//}
	//harbor2.Harbor(i)
}
