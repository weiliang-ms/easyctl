package upgrade

import (
	_ "embed"
	"github.com/spf13/cobra"
)

//go:embed asset/upgrade_kernel.sh
var kernelScript []byte

const kernel = "kernel"

func init() {
	upgradeKernelCmd.Flags().StringVarP(&kernelVersion, "kernel-version", "", "lt", "内核版本 lt|ml")
}

// install kernel
var upgradeKernelCmd = &cobra.Command{
	Use:   "kernel [flags]",
	Short: "upgrade kernel through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		//upgradeKernel()
	},
	Args: cobra.NoArgs,
}

//func upgradeKernel() {
//	if offline && serverListFile == "" {
//		upgradeKernelOffline(filePath)
//	}
//	if offline && serverListFile != "" {
//		list := runner.ParseServerList(serverListFile, runner.CommonServerList{}).Common.Server
//		table.OutputA(list)
//		// todo: 确认交互
//		upgradeKernelOfflineParallel(list)
//	}
//}
//
//// 单机本地离线
//func upgradeKernelOffline(filePath string) {
//	var re runner.ExecResult
//	log.Printf("开始升级安装%s...\n", kernel)
//	re = runner.Shell(fmt.Sprintf("version=%s filepath=%s %s", kernelVersion, filePath, string(kernelScript)))
//	if re.ExitCode != 0 {
//		log.Fatal(re.StdErr)
//	}
//}
//
//func upgradeKernelOfflineParallel(list []runner.Server) {
//
//	var wg sync.WaitGroup
//	ch := make(chan runner.ShellResult, len(list))
//
//	// 拷贝文件
//	dstPath := fmt.Sprintf("/tmp/kernel-%s.tar.gz", kernelVersion)
//	binaryName := "easyctl"
//
//	for _, v := range list {
//		log.Printf("传输数据文件%s至%s:/tmp/kernel-%s.tar.gz...", dstPath, v.Host, kernelVersion)
//		runner.ScpFile(filePath, dstPath, v, 0755)
//		log.Println("-> done 传输完毕...")
//
//		log.Printf("传输数据文件%s至%s:/tmp/%s...", binaryName, v.Host, binaryName)
//		runner.ScpFile(binaryName, fmt.Sprintf("/tmp/%s", binaryName), v, 0755)
//		log.Println("-> done 传输完毕...")
//	}
//
//	cmd := fmt.Sprintf("/tmp/easyctl upgrade kernel --offline-file=%s --offline",
//		fmt.Sprintf("/tmp/kernel-%s.tar.gz", kernelVersion))
//
//	// 并行
//	log.Println("-> 批量安装更新...")
//	for _, v := range list {
//		wg.Add(1)
//		go func(server runner.Server) {
//			defer wg.Done()
//			re := server.RemoteShell(cmd)
//			ch <- re
//		}(v)
//	}
//
//	// ch -> slice
//	var as []runner.ShellResult
//	for target := range ch {
//		as = append(as, target)
//	}
//
//	// 表格输出
//	log.Println("执行结果如下：")
//	table.OutputA(as)
//
//	log.Println("-> 重启主机生效...")
//	wg.Wait()
//}
