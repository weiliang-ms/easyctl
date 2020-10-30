package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/sys"
	"github.com/weiliang-ms/easyctl/util"
)

var CloseServiceForever bool
var closeValidArgs = []string{"firewalld", "selinux", "desktop"}

func init() {

	closeSeLinuxCmd.Flags().BoolVarP(&CloseServiceForever, "forever", "f", false, "Service closed duration.")
	closeFirewalldCmd.Flags().BoolVarP(&CloseServiceForever, "forever", "f", false, "Service closed duration.")

	closeCmd.AddCommand(closeSeLinuxCmd)
	closeCmd.AddCommand(closeFirewalldCmd)
	RootCmd.AddCommand(closeCmd)
}

// close命令
var closeCmd = &cobra.Command{
	Use:   "close [OPTIONS] [flags]",
	Short: "close some service through easyctl",
	Example: "\neasyctl close firewalld" +
		"\neasyctl close firewalld --forever=true" +
		"\neasyctl close selinux" +
		"\neasyctl close selinux --forever=true",
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

// close firewalld命令
var closeFirewalldCmd = &cobra.Command{
	Use:   "firewalld [flags]",
	Short: "close firewalld through easyctl",
	Example: "\neasyctl close firewalld 临时关闭firewalld" +
		"\neasyctl close firewalld --forever=true 永久关闭firewalld" +
		"\neasyctl close firewalld -f 永久关闭firewalld",
	Run: func(cmd *cobra.Command, args []string) {
		closeFirewalld()
	},
	ValidArgs: closeValidArgs,
}

// 关闭selinux
func closeSeLinux() {
	fmt.Printf("#### 关闭selinux服务 ####\n\n")
	sys.CloseSeLinux(CloseServiceForever)
}

// 关闭防火墙

func closeFirewalld() {
	util.PrintTitleMsg("关闭防火墙服务")
	sys.CloseFirewalld(CloseServiceForever)
}
