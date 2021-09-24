package deny

func Firewall(config []byte, debug bool) error {
	return Item(config, debug, disableFirewallShell)
}
