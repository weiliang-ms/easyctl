package resources

const Redhat7RedisServiceContent = "" +
	"[Unit]\n" +
	"Description=Redis persistent key-value database\n" +
	"After=network.target\n" +
	"After=network-online.target\n" +
	"Wants=network-online.target\n\n" +
	"[Service]\n" +
	"ExecStart=/usr/bin/redis-server /etc/redis.conf --supervised systemd\n" +
	"ExecStop=/usr/libexec/redis-shutdown\n" +
	"Type=notify\n" +
	"User=redis\n" +
	"Group=redis\n" +
	"RuntimeDirectory=redis\n" +
	"RuntimeDirectoryMode=0755\n\n" +
	"[Install]\n" +
	"WantedBy=multi-user.target"
