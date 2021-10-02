package deny

import "github.com/sirupsen/logrus"

// Firewall 关闭防火墙
func Firewall(config []byte, logger *logrus.Logger) error {
	return Item(config, logger, disableFirewallShell)
}
