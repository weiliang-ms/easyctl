package set

import "github.com/sirupsen/logrus"

const setTimezoneShell = "ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime"

func Timezone(config []byte, logger *logrus.Logger) error {
	return Config(config, logger, setTimezoneShell)
}
