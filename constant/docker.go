package constant

const (
	Redhat7DockerServiceFilePath = "/etc/systemd/system/docker.service"
)

// /usr/lib/systemd/system/docker.service
const Redhat7DockerServiceContent = "" +
	"[Unit]\nDescription=Docker Application Container Engine\n" +
	"Documentation=https://docs.docker.com\n" +
	"After=network-online.target firewalld.service\n" +
	"Wants=network-online.target\n" +
	"[Service]\n" +
	"Type=notify\n" +
	"# the default is not to use systemd for cgroups because the delegate issues still\n" +
	"# exists and systemd currently does not support the cgroup feature set required\n" +
	"# for containers run by docker\n" +
	"ExecStart=/usr/bin/dockerd\n" +
	"ExecReload=/bin/kill -s HUP \n" +
	"# Having non-zero Limit*s causes performance problems due to accounting overhead\n" +
	"# in the kernel. We recommend using cgroups to do container-local accounting.\n" +
	"LimitNOFILE=infinity\n" +
	"LimitNPROC=infinity\n" +
	"LimitCORE=infinity\n" +
	"# Uncomment TasksMax if your systemd version supports it.\n" +
	"# Only systemd 226 and above support this version.\n" +
	"#TasksMax=infinity\n" +
	"TimeoutStartSec=0\n" +
	"# set delegate yes so that systemd does not reset the cgroups of docker containers\n" +
	"Delegate=yes\n" +
	"# kill only the docker process, not all processes in the cgroup\n" +
	"KillMode=process\n" +
	"# restart the docker process if it exits prematurely\n" +
	"Restart=on-failure\n" +
	"StartLimitBurst=3\n" +
	"StartLimitInterval=60s\n" +
	"[Install]\n" +
	"WantedBy=multi-user.target"
