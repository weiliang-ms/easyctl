package install

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/pkg/runner"
)

// 单机本地离线
func DockerCompose(i runner.Installer) {
	re := runner.ParseServerList(i.ServerListPath, runner.DockerComposeServerList{})
	list := re.Compose.Attribute.Server
	i.Cmd = fmt.Sprintf("\\cp /tmp/docker-compose /usr/bin\nchmod +x /usr/bin/docker-compose")
	i.FileName = "docker-compose"
	if i.Offline && len(list) != 0 {
		offlineRemote(i, list)
	}
}
