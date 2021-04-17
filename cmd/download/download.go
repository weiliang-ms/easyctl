package download

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
	"github.com/weiliang-ms/easyctl/util"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	Url string
)

func init() {

	Cmd.AddCommand(downloadDockerCmd)
	Cmd.AddCommand(downloadHarborCmd)
	Cmd.AddCommand(downloadKernelCmd)
	Cmd.AddCommand(downloadKeepaliveCmd)
	Cmd.AddCommand(downloadDockerComposeCmd)
}

// add命令
var Cmd = &cobra.Command{
	Use:   "download [OPTIONS] [flags]",
	Short: "download soft through easyctl",
	Run: func(cmd *cobra.Command, args []string) {
	},
	Args: cobra.ExactValidArgs(1),
}

func download(url string, name string) {

	log.Printf("下载%s安装介质，地址为：%s\n", name, url)

	var fileName string

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	dataDir := fmt.Sprintf("%s/repo/%s",
		util.CurrentPath(), name)

	log.Printf("创建%s安装介质持久化目录：%s...\n", name, dataDir)

	dirErr := os.MkdirAll(dataDir, 0644)

	if dirErr != nil {
		panic(err)
	}

	arr := strings.Split(url, "/")
	if len(arr) != 0 {
		fileName = arr[len(arr)-1]
	}

	path := fmt.Sprintf("%s/%s", dataDir, fileName)
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	total := res.ContentLength
	reader := io.LimitReader(res.Body, total)

	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(180*time.Millisecond),
	)

	bar := p.Add(total,
		mpb.NewBarFiller("[=>-|"),
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)

	// create proxy reader
	proxyReader := bar.ProxyReader(reader)

	//fmt.Println("开始下载...")
	io.Copy(f, proxyReader)
	p.Wait()

}
