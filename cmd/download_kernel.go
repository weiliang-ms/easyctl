package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	downloadKernelCmd.Flags().StringVarP(&Url, "url", "",
		// todo: 读取配置文件
		"",
		"docker repo resource url")
	downloadKernelCmd.MarkFlagRequired("url")
}

// download docker resources
var downloadKernelCmd = &cobra.Command{
	Use:     "kernel [flags]",
	Short:   "download kernel soft resources through easyctl...",
	Example: "\neasyctl download docker --url=http://192.168.11.222:2222/repo/kernel/kernel-lt.tar.gz",
	Run: func(cmd *cobra.Command, args []string) {
		download(Url, "kernel")
	},
	Args: cobra.NoArgs,
}
