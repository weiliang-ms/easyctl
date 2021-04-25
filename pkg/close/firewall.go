package close

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"log"
)

func (c *Closer) Firewall() {
	c.parseServerList().firewallCmd().closeFirewall()
}

func (c *Closer) firewallCmd() *Closer {
	switch c.Forever {
	case true:
		c.Cmd = disableFirewallShell
	case false:
		c.Cmd = closeFirewallShell
	}

	return c
}

func (c *Closer) closeFirewall() {

	if len(c.ServerList.Server) > 0 {
		for _, v := range c.ServerList.Server {
			log.Printf("关闭%s节点防火墙...", v.Host)
		}
	} else {
		log.Printf("关闭本机防火墙...")
		runner.Shell(c.Cmd)
	}
}
