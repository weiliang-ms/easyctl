package constant

const Redhat7RedisServiceFilePath = "/usr/lib/systemd/system/redis.service"

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

var RedisConfigContent = "" +
	"bind 0.0.0.0\n" +
	"protected-mode yes\n" +
	"port 26379\n" +
	"tcp-backlog 511\n" +
	"timeout 0\n" +
	"tcp-keepalive 300\n" +
	"daemonize yes\n" +
	"supervised no\n" +
	"pidfile /redis/run/redis_26379.pid\n" +
	"loglevel notice\n" +
	"logfile /redis/log/redis_26379.log\n" +
	"databases 16\n" +
	"always-show-logo yes\n" +
	"save 900 1\n" +
	"save 300 10\n" +
	"save 60 10000\n" +
	"stop-writes-on-bgsave-error yes\n" +
	"rdbcompression yes\n" +
	"rdbchecksum yes\n" +
	"dbfilename dump-26379.rdb\n" +
	"dir /redis/lib/\n" +
	"masterauth redis\n" +
	"replica-serve-stale-data yes\n" +
	"replica-read-only yes\n" +
	"repl-diskless-sync no\n" +
	"repl-diskless-sync-delay 5\n" +
	"repl-disable-tcp-nodelay no\n" +
	"replica-priority 100\n" +
	"requirepass redis\n" +
	"lazyfree-lazy-eviction no\n" +
	"lazyfree-lazy-expire no\n" +
	"lazyfree-lazy-server-del no\n" +
	"replica-lazy-flush no\n" +
	"appendonly no\n" +
	"appendfilename \"appendonly.aof\"\n" +
	"appendfsync everysec\n" +
	"no-appendfsync-on-rewrite no\n" +
	"auto-aof-rewrite-percentage 100\n" +
	"auto-aof-rewrite-min-size 64mb\n" +
	"aof-load-truncated yes\n" +
	"aof-use-rdb-preamble yes\n" +
	"lua-time-limit 5000\n" +
	"cluster-enabled yes\n" +
	"cluster-config-file /etc/nodes-26379.conf\n" +
	"slowlog-log-slower-than 10000\n" +
	"slowlog-max-len 128\n" +
	"latency-monitor-threshold 0\n" +
	"notify-keyspace-events \"\"\n" +
	"hash-max-ziplist-entries 512\n" +
	"hash-max-ziplist-value 64\n" +
	"list-max-ziplist-size -2\n" +
	"list-compress-depth 0\n" +
	"set-max-intset-entries 512\n" +
	"zset-max-ziplist-entries 128\n" +
	"zset-max-ziplist-value 64\n" +
	"hll-sparse-max-bytes 3000\n" +
	"stream-node-max-bytes 4096\n" +
	"stream-node-max-entries 100\n" +
	"activerehashing yes\n" +
	"client-output-buffer-limit normal 0 0 0\n" +
	"client-output-buffer-limit replica 256mb 64mb 60\n" +
	"client-output-buffer-limit pubsub 32mb 8mb 60\n" +
	"hz 10\n" +
	"dynamic-hz yes\n" +
	"aof-rewrite-incremental-fsync yes\n" +
	"rdb-save-incremental-fsync yes"
