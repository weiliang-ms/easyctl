# easyctl

基于golang轻量级运维工具集

** 适用平台：** `CentOS6` | `CentOS7`

# 总览

- [安装](#安装)

- [命令]()
  * [add 添加](#add指令集)
    * [user 用户](#创建用户)
  * [close 关闭](#close指令集)
    * [firewalld](#关闭防火墙)
    * [selinux](#关闭selinux)
  * [install 安装](#install指令集)
    * [docker](#安装docker)
    * [nginx](#安装nginx)
    * [redis](#安装redis)
  * [set 设置](#set指令集)
    * [dns 域名解析](#配置dns域名解析)
    * [hosname 主机名](#配置主机名)
    * [timezone 时区](#配置时区)
    * [yum 镜像源](#yum镜像源)
  * [search 查询](#search指令集)
    * [port 端口监听](#端口监听查询)
- [TODO](#todo)
- [开源项目](#开源项目)


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

# 命令介绍

    Usage:
      easyctl [command] [flags]
    
    Available Commands:
      help        Print the version number of easyctl
      search      search something through easyctl
      set         set something through easyctl
      version     Print the version number of easyctl
    
    Flags:
      -h, --help   help for easyctl

# add指令集

## 创建用户

> 添加用户

1.添加可登录的linux用户(password可省，默认密码：user123)

    easyctl add userad -u username -p password
    
2.添加非登录linux用户

    easyctl add -u username --no-login

# close指令集

## 关闭firewalld

> 格式

    easyctl close firewalld [flags]
    
    flags 可选 -f(永久关闭)
    
> 样例

临时关闭firewalld
    
    easyctl close firewalld
    
永久关闭firewalld

    easyctl close firewalld -f

## 关闭selinux

> 格式

    easyctl close selinux [flags]
    
    flags 可选 -f(永久关闭)
    
> 样例

临时关闭selinux
    
    easyctl close selinux
    
永久关闭selinux

    easyctl close selinux -f
    
# install指令集

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
    
## 配置dns域名解析

> 命令格式

    easyctl set dns dns地址

> 使用样例

    easyctl set dns 114.114.114.114
    
## 配置时区

> 使用样例

    easyctl set timezone
    
或

    easyctl set tz
    
默认配置时区为`上海`，暂不支持可选时区

## 配置主机名

> 命令格式

    easyctl set hostname 主机名

> 使用方式

    easyctl set hostname nginx-server1
    
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