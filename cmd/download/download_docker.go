package download

import (
	"github.com/spf13/cobra"
)

func init() {
	downloadDockerCmd.Flags().StringVarP(&Url, "url", "",
		// todo: 读取配置文件
		"",
		"docker repo resource url")
	downloadDockerCmd.MarkFlagRequired("url")
}

// download docker resources
var downloadDockerCmd = &cobra.Command{
	Use:     "docker [flags]",
	Short:   "download linux soft resources through easyctl...",
	Example: "\neasyctl download docker --url=http://192.168.11.222:2222/repo/docker/repo.tar.gz",
	Run: func(cmd *cobra.Command, args []string) {
		download(Url, "docker-ce")
	},
	Args: cobra.NoArgs,
}
