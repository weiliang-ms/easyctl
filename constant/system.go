package constant

import (
	"fmt"
)

func CreateNologinUserCmd(username string) string {
	return fmt.Sprintf("id -u %s &>/dev/null;[ $? -ne 0 ] && [ `id -u` -eq 0 ] && groupadd %s;useradd %s -g %s -s /sbin/nologin -M",
		username, username, username, username)
}

var (
	Redhat6      = "[ \"$(cat /etc/redhat-release|sed -r 's/.* ([0-9]+)\\..*/\\1/')\" == \"6\" ]"
	Redhat7      = "[ \"$(cat /etc/redhat-release|sed -r 's/.* ([0-9]+)\\..*/\\1/')\" == \"7\" ]"
	DeamonReload = "systemctl reload-daemon"
)
