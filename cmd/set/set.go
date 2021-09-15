package set

import (
	"fmt"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"log"
	"sync"
)

var (
	value          string //默认参数
	serverListFile string // 服务器列表配置文件
	configFile     string
	debug          bool
)

// RootCmd set命令
var RootCmd = &cobra.Command{
	Use:   "set [OPTIONS] [flags]",
	Short: "设置指令集",
	Args:  cobra.ExactValidArgs(1),
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "-c", "", "主机列表")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "开启debug模式")

	RootCmd.AddCommand(ulimitCmd)

	RootCmd.AddCommand(dnsCmd)
	RootCmd.AddCommand(hostResolveCmd)
	RootCmd.AddCommand(timeZoneCmd)
	RootCmd.AddCommand(yumRepoCmd)
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

	ch := make(chan runner.ShellResult, len(list.Common.Server))

	// 拷贝文件
	binaryName := "easyctl"

	for _, v := range list.Common.Server {
		log.Printf("传输数据文件%s至%s:/tmp/%s...", binaryName, v.Host, binaryName)
		runner.ScpFile(binaryName, fmt.Sprintf("/tmp/%s", binaryName), v, 0755)
		log.Println("-> done 传输完毕...")
	}

	// 并行
	log.Println("-> 批量配置...")
	for _, v := range list.Common.Server {
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

func level() util.LogLevel {
	level := util.Info
	if debug {
		level = util.Debug
	}
	return level
}
