package close

const denyPingShell = `
sed -i "/net.ipv4.icmp_echo_ignore_all/d" /etc/sysctl.conf
echo "net.ipv4.icmp_echo_ignore_all=1"  >> /etc/sysctl.conf
sysctl -p
`

func (ac *Actuator) Ping() {
	ac.parseServerList().denyPingCmd().execute("禁ping", 0)
}

func (ac *Actuator) denyPingCmd() *Actuator {
	ac.Cmd = denyPingShell
	return ac
}
