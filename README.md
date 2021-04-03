<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [easyctl](#easyctl)
  - [安装](#%E5%AE%89%E8%A3%85)
- [指令集](#%E6%8C%87%E4%BB%A4%E9%9B%86)
  - [add](#add)
    - [user](#user)
  - [close](#close)
    - [firewalld](#firewalld)
    - [selinux](#selinux)
  - [install](#install)
    - [keepalive](#keepalive)
      - [离线](#%E7%A6%BB%E7%BA%BF)
  - [安装docker](#%E5%AE%89%E8%A3%85docker)
  - [安装nginx](#%E5%AE%89%E8%A3%85nginx)
  - [安装redis](#%E5%AE%89%E8%A3%85redis)
- [search指令集](#search%E6%8C%87%E4%BB%A4%E9%9B%86)
  - [端口监听查询](#%E7%AB%AF%E5%8F%A3%E7%9B%91%E5%90%AC%E6%9F%A5%E8%AF%A2)
- [set指令集](#set%E6%8C%87%E4%BB%A4%E9%9B%86)
  - [yum镜像源](#yum%E9%95%9C%E5%83%8F%E6%BA%90)
  - [yum代理配置](#yum%E4%BB%A3%E7%90%86%E9%85%8D%E7%BD%AE)
  - [dns](#dns)
    - [可选参数](#%E5%8F%AF%E9%80%89%E5%8F%82%E6%95%B0)
    - [命令格式](#%E5%91%BD%E4%BB%A4%E6%A0%BC%E5%BC%8F)
    - [使用样例](#%E4%BD%BF%E7%94%A8%E6%A0%B7%E4%BE%8B)
  - [password-less](#password-less)
    - [可选参数](#%E5%8F%AF%E9%80%89%E5%8F%82%E6%95%B0-1)
    - [命令格式](#%E5%91%BD%E4%BB%A4%E6%A0%BC%E5%BC%8F-1)
    - [使用样例](#%E4%BD%BF%E7%94%A8%E6%A0%B7%E4%BE%8B-1)
  - [timezone](#timezone)
    - [可选参数](#%E5%8F%AF%E9%80%89%E5%8F%82%E6%95%B0-2)
    - [命令格式](#%E5%91%BD%E4%BB%A4%E6%A0%BC%E5%BC%8F-2)
    - [使用样例](#%E4%BD%BF%E7%94%A8%E6%A0%B7%E4%BE%8B-2)
  - [配置主机名](#%E9%85%8D%E7%BD%AE%E4%B8%BB%E6%9C%BA%E5%90%8D)
  - [upgrade 命令](#upgrade-%E5%91%BD%E4%BB%A4)
    - [内核](#%E5%86%85%E6%A0%B8)
      - [离线](#%E7%A6%BB%E7%BA%BF-1)
  - [todo](#todo)
  - [开源项目](#%E5%BC%80%E6%BA%90%E9%A1%B9%E7%9B%AE)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# easyctl

基于golang轻量级运维工具集

** 适用平台：** `CentOS7`


## 安装

> 下载上传

[下载release版本](https://github.com/weiliang-ms/easyctl/releases/)

上传至/usr/bin/下

> 添加执行权限

    chmod +x /usr/bin/easyctl
    
> 查看版本信息

    easyctl version
    
> 配置命令补全

    yum install bash-completion -y
    ./easyctl completion bash > /etc/bash_completion.d/easyctl
    source <(./easyctl completion bash)

# 指令集

    Usage:
      easyctl [command] [flags]
    
    Available Commands:
      help        Print the version number of easyctl
      search      search something through easyctl
      set         set something through easyctl
      version     Print the version number of easyctl
    
    Flags:
      -h, --help   help for easyctl

## add

### user

> 添加用户

1.添加可登录的linux用户(password可省，默认密码：user123)

    easyctl add userad -u username -p password
    
2.添加非登录linux用户

    easyctl add -u username --no-login

## close

### firewalld

> 格式

    easyctl close firewalld [flags]
    
    flags 可选 -f(永久关闭)
    
> 样例

临时关闭firewalld
    
    easyctl close firewalld
    
永久关闭firewalld

    easyctl close firewalld -f

### selinux

> 格式

    easyctl close selinux [flags]
    
    flags 可选 -f(永久关闭)
    
> 样例

临时关闭selinux
    
    easyctl close selinux
    
永久关闭selinux

    easyctl close selinux -f
    
## install

### keepalive

安装keepalived

#### 离线

> 1.下载`keepalived`离线仓库

联网主机下执行以下命令:

    sudo docker pull xzxwl/keepalived-repo:latest
    sudo docker run -idt --name keepalived xzxwl/keepalived-repo:latest /bin/bash
    sudo docker cp keepalived:/keepalived.tar.gz ./
    sudo docker rm -f keepalived
    
> 2.安装

初始化生成`server`模板

    ./easyctl init-tmpl keepalived
    
修改`keepalived.yaml`文件内容

    # 虚拟IP
    vip: 192.168.235.150
    # 网卡名称
    interface: ens33
    server:
      - host: 192.168.235.129
        username: root
        password: 1
        port: 22
      - host: 192.168.235.130
        username: root
        password: 1
        port: 22

执行安装

    ./easyctl install keepalived --offline --offline-file=keepalived.tar.gz --server-list=keepalived.yaml
    
安装结果

    ...
    omplete!
    [keepalived] config keepalived...
    [keepalived] boot keepalived...
    Created symlink from /etc/systemd/system/multi-user.target.wants/keepalived.service to /usr/lib/systemd/system/keepalived.service.
    2021/04/02 05:44:56 执行结果如下：
    +-----------------+------------------------------------------------------------------------------------------+------+---------+
    | Host            | Cmd                                                                                      | Code | Status  |
    +-----------------+------------------------------------------------------------------------------------------+------+---------+
    | 192.168.235.129 | /tmp/keepalived.sh ens33 192.168.235.129 192.168.235.130 192.168.235.150 192.168.235.129 | 0    | success |
    | 192.168.235.130 | /tmp/keepalived.sh ens33 192.168.235.129 192.168.235.130 192.168.235.150 192.168.235.130 | 0    | success |
    +-----------------+------------------------------------------------------------------------------------------+------+---------+

## 安装docker

> 格式

    easyctl install docker [flags]
    
    flags 可选 --offline --file=./v19.03.13.tar.gz (离线安装)
    
> 在线安装样例

在线安装`docker`(确保宿主机可访问http://mirrors.aliyun.com)
    
    easyctl install docker
    
> 离线安装样例

**适用于CentOS7**

[下载docker x86压缩包](https://download.docker.com/linux/static/stable/x86_64/)

执行命令安装（--offline --file为必须参数）

    easyctl install docker --file=./docker-19.03.9.tgz --offline

## 安装nginx

> 格式

    easyctl install nginx [flags]
    
    flags 可选 --offline=true --file=./nginx-1.16.0.tar.gz (离线安装)
    
> 样例

在线安装`nginx`(确保宿主机可访问http://mirrors.aliyun.com)
    
    easyctl install nginx
    
## 安装redis

> 格式

    easyctl install redis [flags]
    
flag

    Flags:
      -b, --bind string       Redis bind address (default "0.0.0.0")
      -d, --data string       Redis persistent directory (default "/var/lib/redis")
      -h, --help              help for redis
      -l, --log-file string   Redis logfile directory (default "/var/log/redis")
      -o, --offline           offline mode
      -a, --password string   Redis password (default "redis")
      -p, --port string       Redis listen port (default "6379")
    
> 在线安装样例

在线安装`redis`(确保宿主机可访问http://mirrors.aliyun.com)
    
    easyctl install redis
    
参数定制

    easyctl install redis --bind=192.168.131.36 --data=/var/lib/redis --port=6380 --password=redis567

> 离线安装样例

[下载redis release版本包](http://download.redis.io/releases/),如redis-5.0.5.tar.gz

执行命令安装（其他参数可选，--offline --file为必须参数）

    easyctl install redis --offline --file=./redis-5.0.5.tar.gz

# search指令集

## 端口监听查询

> 命令格式

    easyctl search port 端口值

> 使用样例

    easyctl search port 22

# set指令集

使用方式：easyctl set [options] [flags] 

## yum镜像源


> 配置阿里云yum镜像源

    easyctl set yum --repo=ali
    
或

    easyctl set yum -r=ali
    
> 配置本地镜像源（需手动挂载镜像至/media下：mount -o loop CentOS-7-x86_64-DVD-1908.iso /media）


    easyctl set yum --repo=local
    
或

    easyctl set yum -r=local
 
## yum代理配置

> 配置yum代理

待添加
    
## dns

### 可选参数

    [root@localhost ~]# ./easyctl set dns -h
    easyctl set dns --value
    
    Usage:
      easyctl set dns [flags]
    
    Examples:
    
    easyctl set dns --value=8.8.8.8
    
    Flags:
      -h, --help                 help for dns
          --multi-node           是否配置多节点
          --server-list string   服务器列表 (default "server.yaml")
      -v, --value string         dns 地址...

### 命令格式

> 单节点

    easyctl set dns -v x.x.x.x
    
> 多节点

    easyctl set dns -v x.x.x.x --multi-node

### 使用样例

> 单节点

配置当前主机`dns`

    easyctl set dns -v 8.8.8.8
    
> 多节点

生成`server.yaml`模板文件

    ./easyctl init-tmpl server
    
调整`server.yaml`内容（默认内容如下：）

    server:
      - host: 192.168.235.129
        username: root
        password: 1
        port: 22
      - host: 192.168.235.130
        username: root
        password: 1
        port: 22

配置主机列表内的主机`dns`

    ./easyctl set dns -v 7.7.7.7 --multi-node
  
成功返回  
    
    2021/04/03 02:00:21 <- call back 192.168.235.129
     [dns] 配置成功...
    [dns] 当前dns列表:
    nameserver 114.114.114.114
    nameserver 7.7.7.7
    2021/04/03 02:00:21 <- call back 192.168.235.130
     [dns] 配置成功...
    [dns] 当前dns列表:
    nameserver 114.114.114.114
    nameserver 7.7.7.7
    2021/04/03 02:00:21 执行结果如下：
    +-----------------+----------------+------+---------+
    | Host            | Cmd            | Code | Status  |
    +-----------------+----------------+------+---------+
    | 192.168.235.129 | built-in shell | 0    | success |
    | 192.168.235.130 | built-in shell | 0    | success |
    +-----------------+----------------+------+---------+

## password-less

配置主机间ssh免密登录

### 可选参数

    easyctl set password-less -h
    easyctl set password-less --server-list=xxx
    
    Usage:
      easyctl set password-less [flags]
    
    Examples:
    
    easyctl set password-less --server-list=server.yaml
    
    Flags:
      -h, --help                 help for password-less
          --server-list string   服务器列表 (default "server.yaml")


### 命令格式

    easyctl set password-less

### 使用样例

> 生成`server.yaml`模板文件

    ./easyctl init-tmpl server
    
> 调整`server.yaml`内容（默认内容如下：）

    server:
      - host: 192.168.235.129
        username: root
        password: 1
        port: 22
      - host: 192.168.235.130
        username: root
        password: 1
        port: 22

配置主机列表内的主机间免密访问

    easyctl set password-less
  
成功返回  
    
    
    2021/04/03 15:50:40 生成互信文件
    2021/04/03 15:50:40 -> [192.168.235.129] shell => mkdir -p /root/.ssh
    2021/04/03 15:50:40 <- call back 192.168.235.129
  
    2021/04/03 15:50:40 传输数据文件/root/.ssh/id_rsa至192.168.235.129...
    2021/04/03 15:50:40 -> transfer /root/.ssh/id_rsa to 192.168.235.129
    1.64 KiB / 1.64 KiB [==========================================================| 0s ] 0.00 b/s
    2021/04/03 15:50:40 -> done 传输完毕...
    2021/04/03 15:50:40 传输数据文件/root/.ssh/id_rsa至192.168.235.129...
    2021/04/03 15:50:40 -> transfer /root/.ssh/id_rsa.pub to 192.168.235.129
    408.00 b / 408.00 b [==========================================================| 0s ] 0.00 b/s
    2021/04/03 15:50:40 -> done 传输完毕...
    2021/04/03 15:50:40 传输数据文件/root/.ssh/id_rsa至192.168.235.129...
    2021/04/03 15:50:40 -> transfer /root/.ssh/authorized_keys to 192.168.235.129
    408.00 b / 408.00 b [==========================================================| 0s ] 0.00 b/s
    2021/04/03 15:50:40 -> done 传输完毕...
    2021/04/03 15:50:40 -> [192.168.235.130] shell => mkdir -p /root/.ssh
    2021/04/03 15:50:40 <- call back 192.168.235.130
    
    2021/04/03 15:50:40 传输数据文件/root/.ssh/id_rsa至192.168.235.130...
    2021/04/03 15:50:40 -> transfer /root/.ssh/id_rsa to 192.168.235.130
    1.64 KiB / 1.64 KiB [==========================================================| 0s ] 0.00 b/s
    2021/04/03 15:50:40 -> done 传输完毕...
    2021/04/03 15:50:40 传输数据文件/root/.ssh/id_rsa至192.168.235.130...
    2021/04/03 15:50:40 -> transfer /root/.ssh/id_rsa.pub to 192.168.235.130
    408.00 b / 408.00 b [==========================================================| 0s ] 0.00 b/s
    2021/04/03 15:50:40 -> done 传输完毕...
    2021/04/03 15:50:40 传输数据文件/root/.ssh/id_rsa至192.168.235.130...
    2021/04/03 15:50:40 -> transfer /root/.ssh/authorized_keys to 192.168.235.130
    408.00 b / 408.00 b [==========================================================| 0s ] 0.00 b/s
    2021/04/03 15:50:40 -> done 传输完毕...
    2021/04/03 15:50:40 主机免密配置完毕，请验证...
    
## timezone

默认配置时区为`上海`，暂不支持可选时区

### 可选参数

    easyctl set tz/timezone [value]
    
    Usage:
      easyctl set timezone [flags]
    
    Aliases:
      timezone, tz
    
    Examples:
    
    easyctl set tz/timezone
    
    Flags:
      -h, --help                 help for timezone
          --multi-node           是否配置多节点
          --server-list string   服务器列表 (default "server.yaml")
      -v, --value string         时区 (default "Asia/Shanghai")

### 命令格式

> 单节点

    easyctl set timezone
    
> 多节点

    easyctl set timezone --multi-node

### 使用样例

> 单节点

配置当前主机时区

    easyctl set timezone
    
> 多节点

生成`server.yaml`模板文件

    ./easyctl init-tmpl server
    
调整`server.yaml`内容（默认内容如下：）

    server:
      - host: 192.168.235.129
        username: root
        password: 1
        port: 22
      - host: 192.168.235.130
        username: root
        password: 1
        port: 22

配置主机列表内的主机时区

    ./easyctl set timezone --multi-node
  
成功返回  
    
    2021/04/03 14:51:47 <- call back 192.168.235.129
    
    2021/04/03 14:51:47 <- call back 192.168.235.130
    
    2021/04/03 14:51:47 执行结果如下：
    +-----------------+---------------------------------------------------------+------+---------+
    | Host            | Cmd                                                     | Code | Status  |
    +-----------------+---------------------------------------------------------+------+---------+
    | 192.168.235.129 | \cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime -R | 0    | success |
    | 192.168.235.130 | \cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime -R | 0    | success |
    +-----------------+---------------------------------------------------------+------+---------+


## 配置主机名

> 命令格式

    easyctl set hostname 主机名

> 使用方式

    easyctl set hostname nginx-server1
    
## upgrade 命令

升级`CentOS7`上一些软件

### 内核

更新升级内核

#### 离线

> 1.下载`kernel`离线仓库

联网主机下执行以下命令:

    sudo docker pull xzxwl/kernel-repo:lt
    sudo docker run -idt --name kernel-lt xzxwl/kernel-repo:lt /bin/bash
    sudo docker cp kernel-lt:/data/kernel-lt.tar.gz ./
    sudo docker rm -f kernel-lt
    
> 2.本地更新

    ./easyctl upgrade kernel \
    --offline-file=./kernel-lt.tar.gz --offline
    
> 3.批量更新

初始化生成`server`模板

    ./easyctl init-tmpl server
    
修改`server.yaml`文件内容

    # 默认值
    server:
      - host: 192.168.239.133
        username: root
        password: 123456
        port: 22
      - host: 192.168.239.134
        username: root
        password: 123456
        port: 22

执行安装

    ./easyctl upgrade kernel --offline-file=./kernel-lt.tar.gz --offline --server-list=./server.yaml
    
安装结果

    ...
    2021/04/01 04:53:49 [kernel] check kernel-lt exist ...
    2021/04/01 04:53:49 [kernel] kernel-lt had been installed...
    2021/04/01 04:53:49 0 : CentOS Linux (5.4.108-1.el7.elrepo.x86_64) 7 (Core)
    2021/04/01 04:53:49 1 : CentOS Linux (3.10.0-1062.el7.x86_64) 7 (Core)
    2021/04/01 04:53:49 2 : CentOS Linux (0-rescue-cf09c44eebea4dff8aac64fb57191034) 7 (Core)
    2021/04/01 04:53:49 执行结果如下：
    +-----------------+----------------------------------------------------------------------------+------+---------+
    | Host            | Cmd                                                                        | Code | Status  |
    +-----------------+----------------------------------------------------------------------------+------+---------+
    | 192.168.235.129 | /tmp/easyctl upgrade kernel --offline-file=/tmp/kernel-lt.tar.gz --offline | 0    | success |
    | 192.168.235.130 | /tmp/easyctl upgrade kernel --offline-file=/tmp/kernel-lt.tar.gz --offline | 0    | success |
    +-----------------+----------------------------------------------------------------------------+------+---------+
    2021/04/01 04:53:49 -> 重启主机生效...

    
## todo

1.安全加固脚本（可排除选项）

2.升级软件（在线|离线源码）

3.获取系统信息

4.调整文件描述符|进程数

5.多主机间互信

6.开启端口监听用以测试网络连通性
  
7.关闭某一服务

8.主机host解析

9.添加命令自动补全(已完成)

## 开源项目

- [cobra](https://github.com/spf13/cobra)
- [vssh](https://github.com/yahoo/vssh)