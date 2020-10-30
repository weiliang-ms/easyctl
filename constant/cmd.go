package constant

var (
	LimitOptimizeCmd            = "[ `ulimit -n` -le  65535 ] && echo \"* soft nofile 655350\n* hard nofile 655350\n* soft nproc 65535\n* hard nproc 65535\" >> /etc/security/limits.conf"
	OverCommitMemoryOptimizeCmd = "sed -i '/vm.overcommit_memory = 1/d' /etc/sysctl.conf;echo \"vm.overcommit_memory = 1\" >> /etc/sysctl.conf;sysctl -p"
	RootDetectionCmd            = "[ `id -u` -eq 0 ]"
)
