package install

import (
	"easyctl/asset"
	"easyctl/pkg/runner"
	"fmt"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func init() {
	keepaliveCmd.Flags().BoolVarP(&offline, "offline", "", false, "是否离线安装")
	keepaliveCmd.Flags().StringVarP(&offlineFilePath, "offline-file", "", "keepalived.tar.gz", "离线文件")
	keepaliveCmd.Flags().StringVarP(&serverListFile, "server-list", "", "keepalived.yaml", "服务器批量连接信息")
}

// install keepalive
var keepaliveCmd = &cobra.Command{
	Use:   "keepalived [flags]",
	Short: "install keepalived through easyctl...",
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

	var wg sync.WaitGroup
	re := runner.ParseKeepaliveList(serverListFile)
	list := re.Server

	ch := make(chan runner.ShellResult, len(list))
	// 拷贝文件
	dstPath := "/tmp/keepalived.tar.gz"

	// 生成本地临时文件
	script, _ := asset.Asset("static/script/install_keepalive.sh")
	ioutil.WriteFile("keepalived.sh", script, 0644)

	for _, v := range list {
		log.Printf("传输数据文件%s至%s...", dstPath, v.Host)
		runner.ScpFile(offlineFilePath, dstPath, v, 0755)
		log.Println("-> done 传输完毕...")

		log.Printf("传输数据文件%s至%s:/tmp/%s...", "keepalived.sh", v.Host, "keepalived.sh")
		runner.ScpFile("keepalived.sh", fmt.Sprintf("/tmp/%s", "keepalived.sh"), v, 0755)
		log.Println("-> done 传输完毕...")
	}

	// 清理文件
	os.Remove("keepalived.sh")

	masterIP := list[0].Host
	slaveIP := list[1].Host
	vip := re.Vip
	insterface := re.Interface

	cmd := fmt.Sprintf("/tmp/keepalived.sh %s %s %s %s ",
		insterface, masterIP, slaveIP, vip)

	// 并行
	log.Println("-> 批量安装更新...")
	for _, v := range list {
		wg.Add(1)
		go func(server runner.Server) {
			defer wg.Done()
			re := server.RemoteShell(cmd + server.Host)
			ch <- re
		}(v)
	}
	wg.Wait()
	close(ch)

	// ch -> slice
	var as []runner.ShellResult
	for target := range ch {
		as = append(as, target)
	}

	// 表格输出
	log.Println("执行结果如下：")
	table.OutputA(as)
}
