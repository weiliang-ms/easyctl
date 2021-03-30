package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	downloadDockerComposeCmd.Flags().StringVarP(&Url, "url", "",
		// todo: 读取配置文件
		"",
		"docker repo resource url")
	downloadDockerComposeCmd.MarkFlagRequired("url")
}

// download docker resources
var downloadDockerComposeCmd = &cobra.Command{
	Use:     "docker-compose [flags]",
	Short:   "download docker-compose through easyctl...",
	Example: "\neasyctl download docker --url=http://192.168.11.222:2222/repo/docker-compose",
	Run: func(cmd *cobra.Command, args []string) {
		download(Url, "docker-compose")
	},
	Args: cobra.NoArgs,
}
