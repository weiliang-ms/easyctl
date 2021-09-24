package clean

const cleanDnsShell = `
echo > /etc/resolv.conf
`

func (ac *Actuator) Dns() {
	//ac.parseServerList().cmd().execute("清理dns", 0)
}

func (ac *Actuator) cmd() *Actuator {
	ac.Cmd = cleanDnsShell
	return ac
}
