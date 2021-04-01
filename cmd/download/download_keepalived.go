package download

import (
	"github.com/spf13/cobra"
)

func init() {
	downloadKeepaliveCmd.Flags().StringVarP(&Url, "url", "",
		// todo: 读取配置文件
		"",
		"keepalive repo resource url")
	downloadKernelCmd.MarkFlagRequired("url")
}

// download docker resources
var downloadKeepaliveCmd = &cobra.Command{
	Use:     "keepalive [flags]",
	Short:   "download keepalive soft resources through easyctl...",
	Example: "\neasyctl download keepalive --url=http://192.168.11.222:2222/repo/keepalive/keepalive.tar.gz",
	Run: func(cmd *cobra.Command, args []string) {
		download(Url, "keepalive")
	},
	Args: cobra.NoArgs,
}
