package cmd

import (
	"easyctl/util"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	downloadHarborCmd.Flags().StringVarP(&Url, "url", "",
		// todo: 读取配置文件
		"https://github.com/goharbor/harbor/releases/download/v2.1.4/harbor-offline-installer-v2.1.4.tgz",
		"harbor resource url")
}

// download harbor resources
var downloadHarborCmd = &cobra.Command{
	Use:     "harbor [flags]",
	Short:   "download linux soft resources through easyctl...",
	Example: "\neasyctl download harbor --url=https://github.com/goharbor/harbor/releases/download/v2.1.4/harbor-offline-installer-v2.1.4.tgz",
	Run: func(cmd *cobra.Command, args []string) {
		download(Url, "harbor")
	},
	Args: cobra.NoArgs,
}

func download(url string, name string) {

	fmt.Printf("下载%s安装介质...\n介质下载地址为：%s\n", name, url)

	var fileName string

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	dataDir := fmt.Sprintf("%s/resources/soft/%s",
		util.CurrentPath(), name)

	fmt.Printf("创建%s安装介质持久化目录：%s...\n", name, dataDir)

	dirErr := os.MkdirAll(dataDir, 0644)

	if dirErr != nil {
		panic(err)
	}

	arr := strings.Split(url, "/")
	if len(arr) != 0 {
		fileName = arr[len(arr)-1]
	}

	f, err := os.Create(fmt.Sprintf("%s/%s", dataDir, fileName))
	if err != nil {
		panic(err)
	}

	fmt.Println("开始下载...")
	start := time.Now()
	io.Copy(f, res.Body)
	stop := time.Now()

	fmt.Printf("下载完毕,耗时：%s...\n", stop.Sub(start))
}
