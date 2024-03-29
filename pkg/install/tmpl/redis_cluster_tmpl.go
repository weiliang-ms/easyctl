package tmpl

import (
	"github.com/lithammer/dedent"
	"text/template"
)

var RedisCompileTmpl = template.Must(template.New("compileTmpl").Parse(dedent.Dedent(`
#!/bin/bash
set -e
{{- if .PackageName }}
cd /tmp
if [ ! -f {{ .PackageName }} ];then
  echo /tmp/{{ .PackageName }} Not Found.
  exit 1
fi
tar zxvf {{ .PackageName }}
packageName=$(echo {{ .PackageName }}|sed 's#.tar.gz##g')
echo $packageName
cd $packageName
make -j $(nproc)
make install
{{- end}}
`)))

// RedisConfigTmpl [1]Ports port列表 [2]Password 密码 [3]ClusterEnabled 是否为集群模式
var RedisConfigTmpl = template.Must(template.New("RedisConfigTmpl").Parse(dedent.Dedent(`
{{- if .Ports }}

{{- range .Ports }}
cat > /etc/redis/redis-{{ . }}.conf <<EOF
bind 0.0.0.0
protected-mode yes
port {{ . }}
tcp-backlog 511
timeout 0
tcp-keepalive 300
daemonize yes
supervised no
loglevel notice
databases 16
always-show-logo yes
save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
replica-serve-stale-data yes
replica-read-only yes
repl-diskless-sync no
repl-diskless-sync-delay 5
repl-disable-tcp-nodelay no
replica-priority 100
lazyfree-lazy-eviction no
lazyfree-lazy-expire no
lazyfree-lazy-server-del no
replica-lazy-flush no
appendonly no
appendfilename "appendonly.aof"
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
aof-load-truncated yes
aof-use-rdb-preamble yes
lua-time-limit 5000
slowlog-log-slower-than 10000
slowlog-max-len 128
latency-monitor-threshold 0
notify-keyspace-events ""
hash-max-ziplist-entries 512
hash-max-ziplist-value 64
list-max-ziplist-size -2
list-compress-depth 0
set-max-intset-entries 512
zset-max-ziplist-entries 128
zset-max-ziplist-value 64
hll-sparse-max-bytes 3000
stream-node-max-bytes 4096
stream-node-max-entries 100
activerehashing yes
client-output-buffer-limit normal 0 0 0
client-output-buffer-limit replica 256mb 64mb 60
client-output-buffer-limit pubsub 32mb 8mb 60
hz 10
dynamic-hz yes
aof-rewrite-incremental-fsync yes
rdb-save-incremental-fsync yes
cluster-config-file "/etc/redis/nodes-{{ . }}.conf"
dbfilename "dump-{{ . }}.rdb"
logfile "/var/log/redis-{{ . }}.log"
pidfile "/var/run/redis-{{ . }}.pid"
dir "/var/lib/redis"
EOF
{{- if $.ClusterEnabled }}
cat >> /etc/redis/redis-{{ . }}.conf <<EOF
cluster-enabled yes
EOF
{{- end }}
{{- if $.Password }}
cat >> /etc/redis/redis-{{ . }}.conf <<EOF
requirepass "{{ $.Password }}"
masterauth "{{ $.Password }}"
EOF
{{- end }}
{{- end }}
{{- end }}
`)))

var RedisBootTmpl = template.Must(template.New("redisBootTmpl").Parse(dedent.Dedent(`
{{- if .Ports }}
{{- range .Ports }}
service redis-{{ . }} start
{{- end }}
{{- end }}
`)))

// OpenFirewallPortTmpl todo: 兼容性适配
var OpenFirewallPortTmpl = template.Must(template.New("openFirewallPortTmpl").Parse(dedent.Dedent(`
#!/bin/sh
{{- if .Ports }}
{{- range .Ports }}
firewall-cmd --zone=public --add-port={{ . }}/tcp --permanent || true
{{- end }}
firewall-cmd --reload || true
{{- end }}
`)))

var InitClusterTmpl = template.Must(template.New("InitClusterTmpl").Parse(dedent.Dedent(`
#!/bin/sh
{{- if .EndpointList }}
echo "yes" | /usr/local/bin/redis-cli --cluster create \
{{- range .EndpointList }}
{{ . }} \
{{- end }}
--cluster-replicas 1 {{- if $.Password }} -a {{ $.Password }} {{- end }}
{{- end }}
`)))

var SetRedisServiceTmpl = template.Must(template.New("SetServiceTmpl").Parse(dedent.Dedent(`
{{- if .Ports }}
{{- range .Ports }}
tee /etc/init.d/redis-{{ . }} <<EOF
#!/bin/sh
# chkconfig: 2345 10 90 
# description: Start and Stop redis

#
# Simple Redis init.d script conceived to work on Linux systems
# as it does use of the /proc filesystem.

REDISPORT={{ . }}
EXEC=/usr/local/bin/redis-server #redis-server 路径
CLIEXEC=/usr/local/bin/redis-cli #redis-cli 路径

PIDFILE=/var/run/redis-{{ . }}.pid
CONF="/etc/redis/redis-{{ . }}.conf" #配置地址

case "\$1" in
    start)
        if [ -f \$PIDFILE ]
        then
                echo "\$PIDFILE exists, process is already running or crashed"
        else
                echo "Starting Redis server..."
                \$EXEC \$CONF
        fi
        ;;
    stop)
        if [ ! -f \$PIDFILE ]
        then
                echo "\$PIDFILE does not exist, process is not running"
        else
                PID=\$(cat \$PIDFILE)
                echo "Stopping ..."
                \$CLIEXEC -p \$REDISPORT {{- if $.Password }} -a {{ $.Password }} {{- end }} shutdown
                while [ -x /proc/\${PID} ]
                do
                    echo "Waiting for Redis to shutdown ..."
                    sleep 1
                done
                echo "Redis stopped"
        fi
        ;;
    *)
        echo "Please use start or stop as first argument"
        ;;
esac
EOF

chmod +x /etc/init.d/redis-{{ . }}
chkconfig redis-{{ . }} on
{{- end }}
{{- end }}
`)))