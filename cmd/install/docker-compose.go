package install

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/api/install"
	"github.com/weiliang-ms/easyctl/pkg/runner"
)

// install docker-compose
var dockerComposeCmd = &cobra.Command{
	Use:   "docker-compose [flags]",
	Short: "install docker-compose through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		dockerCompose()
	},
	Args: cobra.NoArgs,
}

// 单机本地离线
func dockerCompose() {
	i := runner.Installer{
		ServerListPath:  serverListFile,
		Offline:         offline,
		OfflineFilePath: offlineFilePath,
	}
	install.DockerCompose(i)
}
