package install

import (
	"easyctl/pkg/run"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	keepaliveCmd.Flags().BoolVarP(&offline, "offline", "", false, "是否离线安装")
	keepaliveCmd.Flags().StringVarP(&offlineFilePath, "offline-file", "", "", "离线文件")
	keepaliveCmd.Flags().StringVarP(&serverListFile, "server-list", "", "./keepalive.yaml", "服务器批量连接信息")
}

// install keepalive
var keepaliveCmd = &cobra.Command{
	Use:   "keepalive [flags]",
	Short: "install keepalive through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		keepalive()
	},
	Args: cobra.NoArgs,
}

func keepalive() {
	if offline {
		keepaliveOffline()
	}
}

func keepaliveOffline() {

	//var wg sync.WaitGroup

	list := run.ParseServerList(serverListFile).Server
	//ch := make(chan run.ShellResult, len(list))

	// 拷贝文件
	dstPath := "/tmp/keepalived.tar.gz"
	binaryName := "easyctl"

	for _, v := range list {
		log.Printf("传输数据文件%s至%s...", dstPath, v.Host)
		run.RemoteWriteFile(offlineFilePath, dstPath, v)
		log.Println("-> done 传输完毕...")

		log.Printf("传输数据文件%s至%s:/tmp/%s...", binaryName, v.Host, binaryName)
		run.RemoteWriteFile(binaryName, fmt.Sprintf("/tmp/%s", binaryName), v)
		log.Println("-> done 传输完毕...")
	}

	//cmd := fmt.Sprintf("/tmp/easyctl upgrade kernel --offline-file=%s --offline",
	//	fmt.Sprintf("/tmp/kernel-%s.tar.gz", kernelVersion))
	//
	//// 并行
	//log.Println("-> 批量安装更新...")
	//for _, v := range list {
	//	wg.Add(1)
	//	go func(server run.Server) {
	//		defer wg.Done()
	//		re := server.RemoteShell(cmd)
	//		ch <- re
	//	}(v)
	//}
	//wg.Wait()
	//close(ch)
	//
	//// ch -> slice
	//var as []run.ShellResult
	//for target := range ch {
	//	as = append(as, target)
	//}
	//
	//// 表格输出
	//log.Println("执行结果如下：")
	//table.OutputA(as)
	//
	//log.Println("-> 重启主机生效...")
}
