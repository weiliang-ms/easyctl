package deny

import "github.com/sirupsen/logrus"

func Firewall(config []byte, logger *logrus.Logger) error {
	return Item(config, logger, disableFirewallShell)
}
