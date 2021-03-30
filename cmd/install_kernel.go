package cmd

import (
	"easyctl/asset"
	"easyctl/util"
	"fmt"
	"github.com/spf13/cobra"
)

var kernelVersion string

const kernel = "kernel"

func init() {
	installKernelCmd.Flags().BoolVarP(&offline, "offline", "", false, "是否离线安装")
	installKernelCmd.Flags().StringVarP(&kernelVersion, "kernel-version", "", "lt", "内核版本 lt|ml")
}

// install kernel
var installKernelCmd = &cobra.Command{
	Use:   "kernel [flags]",
	Short: "install kernel through easyctl...",
	// Example: "\neasyctl download harbor --url=https://github.com/goharbor/harbor/releases/download/v2.1.4/harbor-offline-installer-v2.1.4.tgz",
	Run: func(cmd *cobra.Command, args []string) {
		if offline {
			installKernelOffline()
		} else {

		}
	},
	Args: cobra.NoArgs,
}

// 单机本地离线
func installKernelOffline() {

	script, _ := asset.Asset("script/install_kernel.sh")
	fmt.Printf("开始安装%s...\n", kernel)
	code := util.Run(fmt.Sprintf("version=%s %s", kernelVersion, string(script)))
	if code == 0 {
		fmt.Println("升级内核版本成功，重启生效...")
	} else {
		fmt.Println("升级失败...")
	}

}
