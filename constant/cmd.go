package constant

import "fmt"

var (
	LimitOptimizeCmd            = "if [ `ulimit -n` -le  65535 ];then echo \"* soft nofile 655350\n> > * hard nofile 655350\n> > * soft nproc 65535\n> > * hard nproc 65535\" >> /etc/security/limits.conf\nfi"
	OverCommitMemoryOptimizeCmd = "sed -i '/vm.overcommit_memory = 1/d' /etc/sysctl.conf;echo \"vm.overcommit_memory = 1\" >> /etc/sysctl.conf;sysctl -p"
	RootDetectionCmd            = "[ `id -u` -eq 0 ]"
	EtcRcLocal                  = "/etc/rc.local"
	ChmodX                      = "chmod +x"
	LocalIPCmd                  = "ip addr | awk '/^[0-9]+: / {}; /inet.*global/ {print gensub(/(.*)\\/(.*)/, \"\\\\1\", \"g\", $2)}'|head -1"
	startUnitServiceCmd         = fmt.Sprintf("%s %s", systemctl, start)
	statusUnitServiceCmd        = fmt.Sprintf("%s %s", systemctl, status)
	enableUnitServiceCmd        = fmt.Sprintf("%s %s", systemctl, enable)
	DockerEnableCmd             = fmt.Sprintf("%s %s", enableUnitServiceCmd, Docker)
	DockerStartCmd              = fmt.Sprintf("%s %s", startUnitServiceCmd, Docker)
	DockerStatusCmd             = fmt.Sprintf("%s %s", statusUnitServiceCmd, Docker)
	DisableSelinuxCmd           = "setenforce 0;sed -i \"s#SELINUX=enforcing#SELINUX=disabled#g\" /etc/selinux/config"
	SelinuxStatusCmd            = "sestatus -v|grep \"SELinux status\"|awk '{print $3}'"
	KernelVersionCmd            = "uname -r | grep -o \"^[0-9]\""
	BackupYumRepoCmd            = "mkdir -p /etc/yum.repos.d/`date +%Y%m%d`" + ";" +
		"mv /etc/yum.repos.d/*.repo /etc/yum.repos.d/`date +%Y%m%d` -f"
	CleanYumCacheCmd      = "yum clean all"
	CreateRsaAndRsaPubCmd = "ssh-keygen -t rsa -N '' -f ~/.ssh/id_rsa -q"
)
