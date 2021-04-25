package close

import "github.com/weiliang-ms/easyctl/pkg/runner"

type Closer struct {
	Cmd                string
	Forever            bool
	Remote             bool
	ServerListFilePath string
	ServerList         runner.CommonServerList
}

const (
	closeFirewallShell   = "systemctl disable firewalld --now"
	disableFirewallShell = "systemctl stop firewalld"
	closeSeLinuxShell    = "setenforce 0"
	disableSeLinuxShell  = "setenforce 0 && sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config"
)

func (c *Closer) parseServerList() *Closer {
	c.ServerList = runner.ParseServerList(c.ServerListFilePath, runner.ServerList{}.Common).Common
	return c
}
