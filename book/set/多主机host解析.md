### 配置主机间host解析

采集将多主机间的`hostname`与`IP`解析，过滤`hostname`为`localhost`的条例，配置到`/etc/hosts`中

> 生成默认配置文件

```shell
easyctl set host-resolv
```

> 修改配置文件

`config.yaml`

```yaml
server:
  - host: 10.10.1.[1:3]
    username: root
    password: 111111
    port: 22
excludes:
  - 192.168.235.132
```

> 配置`host`解析

`--debug`输出`debug`日志，可选参数

```shell
easyctl set host-resolv -c config.yaml --debug
```

> 查看解析

```shell
[root@scq-dc01 ~]# cat /etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6

# easyctl hosts BEGIN
10.10.1.1 scq-dc01
10.10.1.2 scq-dc02
10.10.1.3 scq-dc03
# easyctl hosts END
```