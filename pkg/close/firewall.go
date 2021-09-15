package close

func (ac *Actuator) Firewall() {
	ac.parseServerList().firewallCmd().execute("关闭防火墙", 0)
}

func (ac *Actuator) firewallCmd() *Actuator {
	ac.Cmd = disableFirewallShell
	return ac
}
