package set

import (
	"easyctl/pkg/run"
	"flag"
	"fmt"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	"log"
	"sync"
)

const (
	dns        = "dns"
	ali        = "ali"
	idRsa      = "id_rsa"
	ak         = "authorized_keys"
	idRsaPub   = "id_rsa.pub"
	serverList = "server-list"
	isoPath    = "iso-path"
)

var (
	value          string //默认参数
	multiNode      bool   // 是否多节点
	serverListFile string // 服务器列表配置文件

	yumServerListFile   string
	trustServerListFile string // 主机互信主机列表
	imageFilePath       string
	yumRepo             string // 仓库地址
	yumProxy            string // 代理地址
)

// set命令
var RootCmd = &cobra.Command{
	Use:   "set [OPTIONS] [flags]",
	Short: "set something through easyctl",
	//Example: "\neasyctl set dns 114.114.114.114" +
	//"\neasyctl set yum ali" +
	//"\neasyctl set hostname weiliang.com",
	Args: cobra.ExactValidArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cmd.ValidateArgs(setValidArgs)
	},
	ValidArgs: setValidArgs,
}

func init() {
	setYumCmd.Flags().StringVarP(&yumRepo, "repo", "r", "", "Repository address of yum")
	setYumCmd.Flags().StringVarP(&yumProxy, "proxy", "p", "", "Proxy address of yum")
	setYumCmd.Flags().StringVarP(&yumServerListFile, serverList, "", "", "配置yum主机列表")
	setYumCmd.Flags().StringVarP(&imageFilePath, isoPath, "", "", "本机系统版本镜像路径")

	RootCmd.AddCommand(setDNSCmd)
	RootCmd.AddCommand(setPasswordLessCmd)
	RootCmd.AddCommand(setTimeZoneCmd)

	//RootCmd.AddCommand(setYumCmd)
	//RootCmd.AddCommand(setHostnameCmd)
	flag.Parse()

}

// 配置yum子命令
var setYumCmd = &cobra.Command{
	Use:   "yum [flags]",
	Short: "easyctl set yum [flags]",
	Example: "\neasyctl set yum --repo=ali" +
		"\neasyctl set yum --repo=local" +
		"\neasyctl set yum --proxy=http://username:password@192.168.111.222:8080",
	Args: cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		setYum(cmd)
	},
}

func needFlag(cmd *cobra.Command) {
	if cmd.Flags().NFlag() == 0 {
		fmt.Printf("Flags:\n%s", cmd.Flags().FlagUsages())
	}
}

func setYum(cmd *cobra.Command) {
	needFlag(cmd)
	//if yumRepo == local && imageFilePath == "" {
	//	log.Fatal("配置本地yum源，必须通过--iso-path指定iso镜像路径")
	//}
	var yum yum
	yum.setYum()
}

// 单机本地
func local(msg string, cmd string) {
	var re run.ExecResult
	log.Println(msg)
	re = run.Shell(cmd)
	if re.ExitCode != 0 {
		log.Fatal(re.StdErr)
	}
}

// 跨服务器多节点
func multiShell(list run.ServerList, cmd string) {
	var wg sync.WaitGroup

	ch := make(chan run.ShellResult, len(list.Server))

	// 拷贝文件
	binaryName := "easyctl"

	for _, v := range list.Server {
		log.Printf("传输数据文件%s至%s:/tmp/%s...", binaryName, v.Host, binaryName)
		run.RemoteWriteFile(binaryName, fmt.Sprintf("/tmp/%s", binaryName), v, 0755)
		log.Println("-> done 传输完毕...")
	}

	// 并行
	log.Println("-> 批量配置...")
	for _, v := range list.Server {
		wg.Add(1)
		go func(server run.Server) {
			defer wg.Done()
			re := server.RemoteShell(cmd)
			ch <- re
		}(v)
	}
	wg.Wait()
	close(ch)

	// ch -> slice
	var as []run.ShellResult
	for target := range ch {
		as = append(as, target)
	}

	// 表格输出
	log.Println("执行结果如下：")
	table.OutputA(as)
}
