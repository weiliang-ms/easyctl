- [easyctl](#easyctl)
  - [安装](#%E5%AE%89%E8%A3%85)
    - [编译安装最新版](#%E7%BC%96%E8%AF%91%E5%AE%89%E8%A3%85%E6%9C%80%E6%96%B0%E7%89%88)
  - [set指令集](#set%E6%8C%87%E4%BB%A4%E9%9B%86)
    - [配置主机间host解析](#%E9%85%8D%E7%BD%AE%E4%B8%BB%E6%9C%BA%E9%97%B4host%E8%A7%A3%E6%9E%90)
    - [配置主机间免密登录](#%E9%85%8D%E7%BD%AE%E4%B8%BB%E6%9C%BA%E9%97%B4%E5%85%8D%E5%AF%86%E7%99%BB%E5%BD%95)
    - [配置主机文件描述符](#%E9%85%8D%E7%BD%AE%E4%B8%BB%E6%9C%BA%E6%96%87%E4%BB%B6%E6%8F%8F%E8%BF%B0%E7%AC%A6)
    - [配置主机时区](#%E9%85%8D%E7%BD%AE%E4%B8%BB%E6%9C%BA%E6%97%B6%E5%8C%BA)
    - [修改主机root口令](#%E4%BF%AE%E6%94%B9%E4%B8%BB%E6%9C%BAroot%E5%8F%A3%E4%BB%A4)
    
# easyctl

基于`golang`轻量级运维工具集

** 适用平台：** `CentOS7`

## 安装

### 编译安装最新版

```shell
git clone https://github.com/weiliang-ms/easyctl.git
cd easyctl
go build -ldflags "-w -s" -o /usr/local/bin/easyctl
```

## set指令集

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

### 配置主机间免密登录

多主机间配置免密`ssh`登录（基于密钥登录）

> 生成默认配置文件

```shell
easyctl set password-less
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

> 配置免密登录

`--debug`输出`debug`日志，可选参数

```shell
easyctl set password-less -c config.yaml --debug
```

> 测试

`10.10.1.2`为主机列表内的主机

```shell
ssh 10.10.1.2
```

### 配置主机文件描述符

多主机配置文件描述符数量（65535）

> 生成默认配置文件

```shell
easyctl set ulimit
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

> 配置免密登录

`--debug`输出`debug`日志，可选参数

```shell
easyctl set ulimit -c config.yaml --debug
```

### 配置主机时区

多主机配置时区（上海时区）

> 生成默认配置文件

```shell
easyctl set tz
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

> 配置免密登录

`--debug`输出`debug`日志，可选参数

```shell
easyctl set tz -c config.yaml --debug
```

> 测试

```shell
date
```

### 修改主机root口令

> 生成默认配置文件

```shell
easyctl set new-password
```

> 修改配置文件

`config.yaml`

- 调整主机信息，新`root`口令的值

```yaml
server:
  - host: 10.10.10.[1:40]
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
newRootPassword: "3LEPnok84HxYc5"
```

> 运行

`--debug`输出`debug`日志，可选参数

```shell
easyctl set new-password -c config.yaml --debug
```

> 成功样例

```
[root@localhost ~]# ./easyctl set new-password -c config.yaml
I0926 15:14:56.431946  112411 log.go:184] 检测到配置文件中含有IP段，开始解析组装...
I0926 15:14:56.431999  112411 parse.go:113] 解析到IP子网网段为：10.10.1....
I0926 15:14:56.432021  112411 parse.go:117] 解析到IP区间为：1:2...
I0926 15:14:56.432026  112411 parse.go:121] 解析到起始IP为：10.10.1.1...
I0926 15:14:56.432031  112411 parse.go:125] 解析到末尾IP为：10.10.1.2...
I0926 15:14:56.432037  112411 exec.go:43] 开始并行执行命令...
I0926 15:14:56.432084  112411 exec.go:105] [10.10.1.2] 开始执行指令 ->
I0926 15:14:56.432114  112411 exec.go:105] [10.10.1.1] 开始执行指令 ->
I0926 15:14:56.634224  112411 log.go:184] <- 10.10.1.1执行命令成功...
I0926 15:14:56.634472  112411 log.go:184] <- 10.10.1.2执行命令成功...
| IP ADDRESS  |  CMD   | EXIT CODE | RESULT  |        OUTPUT        | EXCEPTION |
|-------------|--------|-----------|---------|----------------------|-----------|
| 10.10.1.1 | ****** |     0     | success | Changing password fo |           |
| 10.10.1.2 | ****** |     0     | success | Changing password fo |           |
```

> 测试

重新连接列表主机