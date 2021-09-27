- [easyctl](#easyctl)
  - [安装](#%E5%AE%89%E8%A3%85)
    - [编译安装最新版](#%E7%BC%96%E8%AF%91%E5%AE%89%E8%A3%85%E6%9C%80%E6%96%B0%E7%89%88)
  - [set指令集](#set%E6%8C%87%E4%BB%A4%E9%9B%86)
    - [配置主机间host解析](#%E9%85%8D%E7%BD%AE%E4%B8%BB%E6%9C%BA%E9%97%B4host%E8%A7%A3%E6%9E%90)
    - [配置主机间免密登录](#%E9%85%8D%E7%BD%AE%E4%B8%BB%E6%9C%BA%E9%97%B4%E5%85%8D%E5%AF%86%E7%99%BB%E5%BD%95)
    - [配置主机文件描述符](#%E9%85%8D%E7%BD%AE%E4%B8%BB%E6%9C%BA%E6%96%87%E4%BB%B6%E6%8F%8F%E8%BF%B0%E7%AC%A6)
    - [配置主机时区](#%E9%85%8D%E7%BD%AE%E4%B8%BB%E6%9C%BA%E6%97%B6%E5%8C%BA)
    - [修改主机root口令](#%E4%BF%AE%E6%94%B9%E4%B8%BB%E6%9C%BAroot%E5%8F%A3%E4%BB%A4)
    - [多主机配置dns](#%E5%A4%9A%E4%B8%BB%E6%9C%BA%E9%85%8D%E7%BD%AEdns)
  - [deny指令集](#deny%E6%8C%87%E4%BB%A4%E9%9B%86)
    - [配置主机禁Ping](#%E9%85%8D%E7%BD%AE%E4%B8%BB%E6%9C%BA%E7%A6%81ping)
    - [配置主机禁用selinux](#%E9%85%8D%E7%BD%AE%E4%B8%BB%E6%9C%BA%E7%A6%81%E7%94%A8selinux)
    - [配置主机禁用防火墙](#%E9%85%8D%E7%BD%AE%E4%B8%BB%E6%9C%BA%E7%A6%81%E7%94%A8%E9%98%B2%E7%81%AB%E5%A2%99)
    
# easyctl

基于`golang`轻量级运维工具集

** 适用平台：** `CentOS7`

## 迭代计划

> 里程碑

- `v0.x.y-alpha`: 添加常用指令集功能（`x`为一级指令集，如`set`;`y`为二级指令集，如`set`指令集中的`dns`子指令集）
- `v1.0.0-beta`: `bug`修复、文档站点、集成`github workflow`
- `v1.0.0-release`: 正式版本

> `v1`已实现功能

- `deny`
  - 防火墙
  - `ping`
  - `selinux`
- `set`
  - 时区
  - `dns`
  - 主机互信
  - 主机`host`解析
  - 文件描述符数
  - 修改`root`口令

> `v2`功能预览

- `windows GUI`

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

### 多主机配置dns

> 生成默认配置文件

```shell
easyctl set dns
```

> 修改配置文件

`config.yaml`

- 调整主机信息
- 调整`dns`地址列表

```yaml
server:
  - host: 10.10.10.[1:40]
    username: root
    privateKeyPath: "" # ~/.ssh/id_rsa，为空默认走password登录；不为空默认走密钥登录
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
dns:
  - 114.114.114.114
  - 8.8.8.8
```

> 运行

`--debug`输出`debug`日志，可选参数

```shell
easyctl set dns -c config.yaml --debug
```

> 测试

任意主机列表内的主机执行：

```shell
cat /etc/hosts
```

## deny指令集

### 配置主机禁Ping

配置主机禁`Ping`

> 生成默认配置文件

```shell
easyctl deny ping
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

> 配置

`--debug`输出`debug`日志，可选参数

```shell
easyctl deny ping -c config.yaml --debug
```

### 配置主机禁用selinux

> 生成默认配置文件

```shell
easyctl deny selinux
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

> 配置

`--debug`输出`debug`日志，可选参数

```shell
easyctl deny selinux -c config.yaml --debug
```

### 配置主机禁用防火墙

> 生成默认配置文件

```shell
easyctl deny firewall
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

> 配置

`--debug`输出`debug`日志，可选参数

```shell
easyctl deny firewall -c config.yaml --debug
```
## export指令集

### chart导出

从`harbor`中批量下载`chart`文件

> 背景说明

通常`chart`应用存放于`harbor`进行管理，但`harbor`（v2.1.4）只支持单个`chart`文件下载。

对于离线批量分发场景，比较难受，特开发该功能。

> 生成默认配置文件

```shell
[root@localhost ~]# easyctl export chart
I0927 09:30:47.685246   14817 export.go:44] 检测到配置文件参数为空，生成配置文件样例 -> config.yaml
```

> 修改配置文件

`vi config.yaml`

```yaml
helm-repo:
  endpoint: 10.10.1.3:80   # harbor访问地址
  domain: harbor.wl.io      # harbor域
  username: admin           # harbor用户
  password: 123456          # harbor密码
  preserveDir: /root/charts # chart包持久化目录
  package: true             # 是否打成tar包
  repo-name: charts         # chart repo harbor内的名称
```

> 配置

`--debug`输出`debug`日志，可选参数

```shell
[root@node1 ~]# easyctl export chart -c config.yaml
INFO[0000] 解析chart仓库配置...
INFO[0000] 待导出chart数量为: 135
INFO[0000] 导出chart...
INFO[0000] 创建目录: /root/charts
INFO[0000] 逐一导出chart中...
INFO[0002] 导出完毕，chart总数为:135
```