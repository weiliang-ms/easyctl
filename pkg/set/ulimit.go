package set

import (
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

const ulimitShell = `
sed -i ':a;$!{N;ba};s@# easyctl ulimit BEGIN.*# easyctl ulimit END@@' /etc/security/limits.conf
sed -i '/^$/N;/\n$/N;//D' /etc/security/limits.conf

cat >> /etc/security/limits.conf <<EOF
# easyctl ulimit BEGIN
* soft nofile 65536
* hard nofile 65536
* soft nproc 65536
* hard nproc 65536
# easyctl ulimit END
EOF
`

// Ulimit 设置文件描述符
func Ulimit(item command.OperationItem) error {

	return Config(item.B, item.Logger, ulimitShell)
}
