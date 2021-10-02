package deny

import "github.com/sirupsen/logrus"

// Ping 禁ping
func Ping(config []byte, logger *logrus.Logger) error {
	return Item(config, logger, denyPingShell)
}
