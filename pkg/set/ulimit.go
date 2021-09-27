package set

import (
	_ "embed"
	"github.com/sirupsen/logrus"
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

func Ulimit(b []byte, logger *logrus.Logger) error {

	return Config(b, logger, ulimitShell)
}
