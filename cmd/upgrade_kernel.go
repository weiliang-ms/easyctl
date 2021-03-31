package cmd

import (
	"easyctl/asset"
	"easyctl/shell"
	"fmt"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var (
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
		list := parseServerList(serverListFile).Server
		table.OutputA(list)
		// todo: 确认交互
		upgradeKernelOfflineParallel(list)
	}
}

// 单机本地离线
func upgradeKernelOffline(filePath string) {
	var re shell.ExecResult
	script, _ := asset.Asset("static/script/upgrade_kernel.sh")
	log.Printf("开始升级安装%s...\n", kernel)
	re = shell.Run(fmt.Sprintf("version=%s filepath=%s %s", kernelVersion, filePath, string(script)))
	if re.ExitCode != 0 {
		log.Fatal(re.StdErr)
	}
}

func upgradeKernelOfflineParallel(list []Server) {

	// 拷贝文件
	filePath := fmt.Sprintf("/tmp/kernel-%s.tar.gz", kernelVersion)
	b, err := ioutil.ReadFile(offlineFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, v := range list {
		log.Printf("传输数据文件%s至%s...", filePath, v.Host)
		remoteWriteFile(filePath, b, v)
		log.Println("传输完毕...")
	}
}
