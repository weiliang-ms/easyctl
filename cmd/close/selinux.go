package close

import (
	"easyctl/pkg/runner"
	"github.com/spf13/cobra"
	"log"
)

func init()  {
	closeSeLinuxCmd.Flags().BoolVarP(&remote,"remote","",false,"是否关闭远程主机seLinux")
	closeSeLinuxCmd.Flags().StringVarP(&serverListFile, "server-list", "", "server.yaml", "服务器列表连接信息")

}

// 关闭防火墙
var closeSeLinuxCmd = &cobra.Command{
	Use:     "selinux",
	Short:   "easyctl close selinux",
	Run: func(cmd *cobra.Command, args []string) {
		closeSeLinux()
	},
}

// 关闭防火墙
func closeSeLinux()  {
	var list []runner.Server
	if remote {
		list = runner.ParseServerList(serverListFile).Server
	}
	cmd := "setenforce 0 && sed -i \"s#SELINUX=enforcing#SELINUX=disabled#g\" /etc/selinux/config"
	close(cmd,list)
	log.Println("重启生效...")
}