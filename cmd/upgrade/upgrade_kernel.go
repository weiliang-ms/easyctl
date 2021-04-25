package upgrade

import (
	"fmt"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/asset"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"log"
	"sync"
)

var (
	offline         bool
	kernelVersion   string
	offlineFilePath string
	serverListFile  string
)

const kernel = "kernel"

func init() {
	upgradeKernelCmd.Flags().BoolVarP(&offline, "offline", "", false, "是否离线安装")
	upgradeKernelCmd.Flags().StringVarP(&offlineFilePath, "offline-file", "", "", "离线文件")
	upgradeKernelCmd.Flags().StringVarP(&serverListFile, "server-list", "", "", "服务器批量连接信息")
	upgradeKernelCmd.Flags().StringVarP(&kernelVersion, "kernel-version", "", "lt", "内核版本 lt|ml")
	//upgradeKernelCmd.MarkFlagRequired("kernel-version")
}

// install kernel
var upgradeKernelCmd = &cobra.Command{
	Use:   "kernel [flags]",
	Short: "upgrade kernel through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		upgradeKernel()
	},
	Args: cobra.NoArgs,
}

func upgradeKernel() {
	if offline && serverListFile == "" {
		upgradeKernelOffline(offlineFilePath)
	}
	if offline && serverListFile != "" {
		list := runner.ParseServerList(serverListFile, runner.CommonServerList{}).Common.Server
		table.OutputA(list)
		// todo: 确认交互
		upgradeKernelOfflineParallel(list)
	}
}

// 单机本地离线
func upgradeKernelOffline(filePath string) {
	var re runner.ExecResult
	script, _ := asset.Asset("static/script/upgrade_kernel.sh")
	log.Printf("开始升级安装%s...\n", kernel)
	re = runner.Shell(fmt.Sprintf("version=%s filepath=%s %s", kernelVersion, filePath, string(script)))
	if re.ExitCode != 0 {
		log.Fatal(re.StdErr)
	}
}

func upgradeKernelOfflineParallel(list []runner.Server) {

	var wg sync.WaitGroup
	ch := make(chan runner.ShellResult, len(list))

	// 拷贝文件
	dstPath := fmt.Sprintf("/tmp/kernel-%s.tar.gz", kernelVersion)
	binaryName := "easyctl"

	for _, v := range list {
		log.Printf("传输数据文件%s至%s:/tmp/kernel-%s.tar.gz...", dstPath, v.Host, kernelVersion)
		runner.ScpFile(offlineFilePath, dstPath, v, 0755)
		log.Println("-> done 传输完毕...")

		log.Printf("传输数据文件%s至%s:/tmp/%s...", binaryName, v.Host, binaryName)
		runner.ScpFile(binaryName, fmt.Sprintf("/tmp/%s", binaryName), v, 0755)
		log.Println("-> done 传输完毕...")
	}

	cmd := fmt.Sprintf("/tmp/easyctl upgrade kernel --offline-file=%s --offline",
		fmt.Sprintf("/tmp/kernel-%s.tar.gz", kernelVersion))

	// 并行
	log.Println("-> 批量安装更新...")
	for _, v := range list {
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

	log.Println("-> 重启主机生效...")
}
