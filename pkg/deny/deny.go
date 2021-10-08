package deny

const (
	disableFirewallShell = "systemctl disable firewalld --now"
	denyPingShell        = `
sed -i "/net.ipv4.icmp_echo_ignore_all/d" /etc/sysctl.conf
echo "net.ipv4.icmp_echo_ignore_all=1"  >> /etc/sysctl.conf
sysctl -p
`
	closeSELinuxShell = `
if [ "$(getenforce)" == "Disabled" ];then
	echo "已关闭，无需重复关闭"
	exit 0
fi
setenforce 0
sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config
`
)
