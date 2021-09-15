package install

import (
	_ "embed"
	"fmt"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

//go:embed asset/install-keepalive-script
var keepaliveScript []byte

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
	re := runner.ParseServerList(serverListFile, runner.KeepaliveServerList{})
	list := re.Keepalive.Attribute.Server

	ch := make(chan runner.ShellResult, len(list))
	// 拷贝文件
	dstPath := "/tmp/keepalived.tar.gz"

	// 生成本地临时文件
	//script, _ := asset.Asset("static/script/install_keepalive.sh")
	ioutil.WriteFile("keepalived.sh", keepaliveScript, 0644)

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

	vip := re.Keepalive.Attribute.Vip
	insterface := re.Keepalive.Attribute.Interface

	cmd := fmt.Sprintf("interface_name=%s virtual_ip=%s /tmp/keepalived.sh",
		insterface, vip)

	// 并行
	log.Println("-> 批量安装更新...")
	for _, v := range list {
		wg.Add(1)
		go func(server runner.Server) {
			defer wg.Done()
			r, peer := role(list, server.Host)
			re := server.RemoteShell(fmt.Sprintf("role=%s peer_ip=%s local_ip=%s %s", r, peer, server.Host, cmd))
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

func role(server []runner.Server, ip string) (role string, peerIP string) {
	if len(server) != 2 {
		log.Fatal("error settings...")
	}
	if server[0].Host == ip {
		return "MASTER", server[1].Host
	} else {
		return "BACKUP", server[0].Host
	}
}
