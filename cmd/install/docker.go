package install

import (
	"github.com/spf13/cobra"
)

// install docker
var dockerCmd = &cobra.Command{
	Use:   "docker-ce [flags]",
	Short: "install docker through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		docker()
	},
}

var dataDir string

func init() {
	dockerCmd.Flags().StringVarP(&dataDir, "data-dir", "", "/var/lib/docker", "数据存储目录")
}

// 单机本地离线
func docker() {
	//i := runner.Installer{
	//	DataDir:         dataDir,
	//	ServerListPath:  serverListFile,
	//	Offline:         offline,
	//	OfflineFilePath: offlineFilePath,
	//}
	//install.Docker(i)
}
