package deny

import "github.com/sirupsen/logrus"

func Ping(config []byte, logger *logrus.Logger) error {
	return Item(config, logger, denyPingShell)
}
