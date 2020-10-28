package constant

import (
	"fmt"
)

func CreateNologinUserCmd(username string) string {
	return fmt.Sprintf("id -u %s && groupadd %s;useradd %s -g %s -s /sbin/nologin -M",
		username, username, username, username)
}
