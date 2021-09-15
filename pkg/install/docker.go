package install

import (
	"github.com/pkg/errors"
	"github.com/weiliang-ms/easyctl/pkg/util/manager"
	"github.com/weiliang-ms/easyctl/pkg/util/ssh"
	"strings"
)

// 单机本地离线
func Docker(mgr *manager.Manager) error {

	createTasks := []manager.Task{
		{Task: parseServerList, ErrMsg: "Failed to parse server list which server will install docker..."},
		{Task: parseMediaType, ErrMsg: "Failed to parse which method to install docker..."},
		{Task: install, ErrMsg: "Failed to parse which method to install docker..."},
	}

	for _, step := range createTasks {
		if err := step.Run(mgr); err != nil {
			return errors.Wrap(err, step.ErrMsg)
		}
	}

	return nil
}

func parseServerList(mgr *manager.Manager) error {
	err, list := ssh.ParseServerList(mgr.ServerListFile, ssh.DockerServerList{})
	if err != nil {
		return errors.Wrap(err, "Failed to find kube offline binaries...")
	}
	mgr.Servers = list.Docker.Attribute.Servers
	return nil
}

// 解析安装介质类型
func parseMediaType(mgr *manager.Manager) error {

	if mgr.OfflineFile == "" {
		mgr.MediaType = manager.YUM
	}

	if strings.HasSuffix(mgr.OfflineFile, ".tgz") {
		mgr.MediaType = manager.BINARY
	} else if strings.HasSuffix(mgr.OfflineFile, ".tar.gz") {
		mgr.MediaType = manager.RPM
	}

	mgr.Logger.Printf("docker安装方式为: %s", mgr.MediaType)
	return nil
}

// 本地安装
func install(mgr *manager.Manager) error {
	mgr.RunTaskOnNodes(mgr.Servers, installDockerOnNode, true)
	return nil
}

func installDockerOnNode(mgr *manager.Manager, nodes *ssh.Server) error {
	return nil
}
