# easyctl

基于golang轻量级运维工具集

** 适用平台：** `CentOS6` | `CentOS7`

> 使用方式

1.[下载release版本]()

2.上传至/usr/sbin/下

3.添加执行权限

    chmod +x /usr/sbin/nginx
    
4.查看版本信息

    easyctl version
    
> 版本说明

    vx.y.z

x为主版本号

y为奇数：alpha测试版本

y为偶数：beta稳定版

z为y的补丁号

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

## add命令

> 添加用户

1.添加可登录的linux用户(password可省，默认密码：user123)

    easyctl add username password
    
2.添加非登录linux用户

    easyctl add username --no-login=true

## set命令

使用方式：easyctl set [options] [flags] 

### yum配置 

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
    
### DNS
    
配置DNS地址

> 配置dns

使用方式

    easyctl set dns 114.114.114.114

### hostname

配置hostname

> 配置hostname

使用方式

    easyctl set hostname nginx-server1
    
## TODO

1.安全加固脚本（可排除选项）

2.升级软件（在线|离线源码）

3.获取系统信息

4.调整文件描述符|进程数

5.多主机间互信

6.开启端口监听用以测试网络连通性
  
