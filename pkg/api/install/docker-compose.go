package install

//
//import (
//	"fmt"
//	"github.com/weiliang-ms/easyctl/pkg/runner"
//	"log"
//)
//
//// 单机本地离线
//func DockerCompose(i runner.Installer) {
//	var list []runner.Server
//
//	if i.ServerListPath != "" {
//		re := runner.ParseServerList(i.ServerListPath, runner.DockerComposeServerList{})
//		list = re.Compose.Attribute.Server
//	}
//
//	i.Cmd = fmt.Sprintf("sudo \\cp %s /usr/bin && sudo chmod +x /usr/bin/docker-compose", i.OfflineFilePath)
//
//	if i.Offline && len(list) != 0 {
//		offlineRemote(i, list)
//	} else if i.Offline && len(list) == 0 {
//		log.Println("本地安装docker-compose...")
//		runner.Shell(i.Cmd)
//		runner.Shell("docker-compose version")
//	}
//}
