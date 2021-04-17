package install

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/asset"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"log"
	"sync"
)

// 单机本地离线
func Docker(i runner.Installer) {
	re := runner.ParseDockerServerList(i.ServerListPath)
	list := re.Docker.Server
	script, _ := asset.Asset("static/script/install_offline_docker.sh")
	i.Cmd = fmt.Sprintf(string(script))
	i.FileName = "docker-ce.tar.gz"

	if i.Offline && len(list) != 0 {
		offlineRemote(i, list)
	}
}

func offlineRemote(i runner.Installer, list []runner.Server) {
	var wg sync.WaitGroup

	ch := make(chan runner.ShellResult, len(list))
	// 拷贝文件
	dstPath := fmt.Sprintf("/tmp/%s", i.FileName)

	// 生成本地临时文件
	for _, v := range list {
		runner.ScpFile(i.OfflineFilePath, dstPath, v, 0755)
		log.Println("-> transfer done ...")
	}

	// 并行
	log.Println("-> 批量安装...")
	for _, v := range list {
		wg.Add(1)
		go func(server runner.Server) {
			defer wg.Done()
			re := server.RemoteShell(i.Cmd)
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
}

func boot(cmd string, list []runner.Server) {
	// 生成本地临时文件
	for _, v := range list {
		v.RemoteShell(cmd)
		log.Println("-> boot service ...")
	}
}
