package cmd

import (
	"easyctl/sys"
	"fmt"
	"github.com/spf13/cobra"
)

var CloseServiceForever bool
var closeValidArgs = []string{"firewalld", "selinux", "desktop"}

func init() {

	closeSeLinuxCmd.Flags().BoolVarP(&CloseServiceForever, "forever", "f", false, "Service closed duration.")
	closeCmd.AddCommand(closeSeLinuxCmd)
	rootCmd.AddCommand(closeCmd)
}

// close命令
var closeCmd = &cobra.Command{
	Use:   "close [OPTIONS] [flags]",
	Short: "close some service through easyctl",
	Example: "\neasyctl close firewalld" +
		"\neasyctl close firewalld --forever=true",
	Run: func(cmd *cobra.Command, args []string) {
	},
	ValidArgs: closeValidArgs,
	Args:      cobra.ExactValidArgs(1),
}

// close selinux命令
var closeSeLinuxCmd = &cobra.Command{
	Use:   "selinux [flags]",
	Short: "close selinux through easyctl",
	Example: "\neasyctl close selinux 暂时关闭selinux" +
		"\neasyctl close selinux --forever=true 永久关闭selinux",
	Run: func(cmd *cobra.Command, args []string) {
		closeSeLinux()
	},
	ValidArgs: closeValidArgs,
}

// 关闭selinux
func closeSeLinux() {
	fmt.Printf("#### 关闭selinux服务 ####\n\n")
	sys.CloseSeLinux(CloseServiceForever)
}
