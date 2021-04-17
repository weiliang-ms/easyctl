package set

import (
	"flag"
	"fmt"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"log"
	"sync"
)

var (
	value          string //默认参数
	multiNode      bool   // 是否多节点
	serverListFile string // 服务器列表配置文件
	aliRepo        bool   // 阿里云镜像源yum仓库
	tsinghuaRepo   bool   // 清华镜像源yum仓库
)

// set命令
var RootCmd = &cobra.Command{
	Use:   "set [OPTIONS] [flags]",
	Short: "set something through easyctl",
	Args:  cobra.ExactValidArgs(1),
}

func init() {

	RootCmd.AddCommand(setDNSCmd)
	RootCmd.AddCommand(setPasswordLessCmd)
	RootCmd.AddCommand(setTimeZoneCmd)
	RootCmd.AddCommand(setYumRepoCmd)

	//RootCmd.AddCommand(setHostnameCmd)
	flag.Parse()

}

// 单机本地
func local(msg string, cmd string) {
	var re runner.ExecResult
	log.Println(msg)
	re = runner.Shell(cmd)
	if re.ExitCode != 0 {
		log.Fatal(re.StdErr)
	}
}

// 跨服务器多节点
func multiShell(list runner.ServerList, cmd string) {
	var wg sync.WaitGroup

	ch := make(chan runner.ShellResult, len(list.Server))

	// 拷贝文件
	binaryName := "easyctl"

	for _, v := range list.Server {
		log.Printf("传输数据文件%s至%s:/tmp/%s...", binaryName, v.Host, binaryName)
		runner.ScpFile(binaryName, fmt.Sprintf("/tmp/%s", binaryName), v, 0755)
		log.Println("-> done 传输完毕...")
	}

	// 并行
	log.Println("-> 批量配置...")
	for _, v := range list.Server {
		wg.Add(1)
		go func(server runner.Server) {
			defer wg.Done()
			re := server.RemoteShell(cmd)
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
