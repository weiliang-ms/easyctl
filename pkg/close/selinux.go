package close

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"log"
)

func (c *Closer) SeLinux() {
	c.parseServerList().seLinuxCmd().closeSeLinux()
}

func (c *Closer) seLinuxCmd() *Closer {
	switch c.Forever {
	case true:
		c.Cmd = disableSeLinuxShell
	case false:
		c.Cmd = closeSeLinuxShell
	}

	return c
}

func (c *Closer) closeSeLinux() {

	if len(c.ServerList.Server) > 0 {
		for _, v := range c.ServerList.Server {
			log.Printf("关闭%s节点selinux, 重启生效...", v.Host)
		}
	} else {
		log.Printf("关闭本机selinux, 重启生效...")
		runner.Shell(c.Cmd)
	}
}
