# easyctl

基于golang轻量级运维工具集

** 适用平台：** `CentOS6` | `CentOS7`

# 总览

- [安装](#安装)

- [命令]()
  * [add 添加](#add指令集)
    * [user 用户](#创建用户)
  * [set 设置](#set指令集)
    * [dns 域名解析](#配置dns域名解析)
    * [hosname 主机名](#配置主机名)
    * [timezone 时区](#配置时区)
    * [yum 镜像源](#yum镜像源)
  * [search 查询](#search)
- [TODO](#todo)


## 安装

1.[下载release版本](https://github.com/weiliang-ms/easyctl/releases/)

2.上传至/usr/sbin/下

3.添加执行权限

    chmod +x /usr/sbin/easyctl
    
4.查看版本信息

    easyctl version

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

    easyctl add username password
    
2.添加非登录linux用户

    easyctl add username --no-login=true

# set指令集

使用方式：easyctl set [options] [flags] 

## yum镜像源

> 配置yum镜像源

a.配置阿里云yum镜像源

    easyctl set yum --repo=ali
    
或

    easyctl set yum -r=ali
    
b.配置本地镜像源（需手动挂载镜像至/media下：mount -o loop CentOS-7-x86_64-DVD-1908.iso /media）


    easyctl set yum --repo=local
    
或

    easyctl set yum -r=local
 
> 配置yum代理

待添加
    
## 配置dns域名解析
    
配置DNS地址

> 配置dns

使用方式

    easyctl set dns 114.114.114.114
    
## 配置时区

> 配置时区

使用方式

    easyctl set timezone
    
或

    easyctl set tz
    
默认配置时区为`上海`，暂不支持可选时区

## 配置主机名

配置hostname

> 配置hostname

使用方式

    easyctl set hostname nginx-server1
    
## todo

1.安全加固脚本（可排除选项）

2.升级软件（在线|离线源码）

3.获取系统信息

4.调整文件描述符|进程数

5.多主机间互信

6.开启端口监听用以测试网络连通性
  
7.关闭某一服务
