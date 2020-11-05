package constant

var (
	LimitOptimizeCmd            = "if [ `ulimit -n` -le  65535 ];then echo \"* soft nofile 655350\n> > * hard nofile 655350\n> > * soft nproc 65535\n> > * hard nproc 65535\" >> /etc/security/limits.conf\nfi"
	OverCommitMemoryOptimizeCmd = "sed -i '/vm.overcommit_memory = 1/d' /etc/sysctl.conf;echo \"vm.overcommit_memory = 1\" >> /etc/sysctl.conf;sysctl -p"
	RootDetectionCmd            = "[ `id -u` -eq 0 ]"
	EtcRcLocal                  = "/etc/rc.local"
	ChmodX                      = "chmod +x"
	LocalIPCmd                  = "ip a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $2}'|tr -d \"addr:\"|awk '{sub(/.{3}$/,\"\")}1'"
)
