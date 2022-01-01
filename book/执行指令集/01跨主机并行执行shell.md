## 跨主机并行执行shell

### 版本&兼容性

> 版本支持

- [v0.7.16-alpha以上](https://github.com/weiliang-ms/easyctl/releases/)

> 兼容性

- [x] `CentOS6`
- [x] `CentOS7`

### 使用方式

> 参考以下链接进行安装

- [安装说明文档](../-安装文档/README.md)

> 生成默认配置文件

```shell
$ easyctl exec shell
INFO[0000] 生成配置文件样例, 请携带 -c 参数重新执行 -> config.yaml
```

> 修改配置文件

`config.yaml`, 修改主机列表。`easyctl`根据主机列表`ssh`远程至目标主机执行`shell`

```yaml
server:
  - host: 10.10.10.[1:3]
    username: root
    privateKeyPath: ~/.ssh/id_rsa
    password: ""
    port: 22
excludes:
  - 192.168.235.132
script: "1.sh"
```

> 执行

添加`--debug`可以输出详细内容。

```shell
$ easyctl exec shell -c config.yaml --debug
```

### 配置项说明

- 主机配置段：该段配置远程执行`shell`的主机信息，字段说明如下
    - `host: 10.10.10.[1:3]` 主机地址段，适用于`ip`连续场景。分隔符可以为`[1:3]`、`1-2`、`[1-2]`、`1:2`
    - `username`: 远程主机`ssh`用户名称，缺省值为`root`
    - `password`: 对应`username`的密码
    - `privateKeyPath`: `ssh`私钥路径
    - `port`: `ssh`端口，默认`22`
    - `excludes`: 排除`host`地址段内的`ip`地址列表

`privateKeyPath`优先级高于`password`:

1. `privateKeyPath`为空，取`password`值，`ssh`使用密码登录方式
2. `privateKeyPath`非空，取`privateKeyPath`值，`ssh`使用密钥登录方式

```yaml
server:
  - host: 10.10.10.[1:3]
    username: root
    privateKeyPath: ~/.ssh/id_rsa
    password: ""
    port: 22
excludes:
  - 192.168.235.132
```

- 脚本配置:
    - `script: "date"`: 远程执行的`shell`指令，适用于运行单个`shell`指令场景
    - `script: "1.sh"`: 远程执行的`shell`脚本，适用于运行多个`shell`指令场景

### 配置样例

> 1.主机: `10.10.10.1-10.10.10.10`执行`date`指令，使用密钥登录方式

```yaml
server:
  - host: 10.10.10.[1:10]
    username: root
    privateKeyPath: ~/.ssh/id_rsa
    password: ""
    port: 22
excludes:
  - 192.168.235.132
script: "date"
```

> 2.主机: `10.10.10.1-10.10.10.10`执行`date`指令，使用密码登录方式

```yaml
server:
  - host: 10.10.10.[1:10]
    username: root
    # privateKeyPath: ~/.ssh/id_rsa
    password: "123456"
    port: 22
excludes:
  - 192.168.235.132
script: "date"
```

> 3.主机: `10.10.10.1`、`10.10.10.3`、`10.10.10.4`执行`date`指令，使用密码登录方式

```yaml
server:
  - host: 10.10.10.[1:4]
    username: root
    # privateKeyPath: ~/.ssh/id_rsa
    password: "123456"
    port: 22
excludes:
  - 10.10.10.2
script: "date"
```

> 4.主机: `10.10.10.1`、`10.10.10.3`、`10.10.10.4`执行`shell`脚本，使用密码登录方式

```yaml
server:
  - host: 10.10.10.[1:4]
    username: root
    # privateKeyPath: ~/.ssh/id_rsa
    password: "123456"
    port: 22
excludes:
  - 10.10.10.2
script: "/root/modify-sysctl.sh"
```

`/root/modify-sysctl.sh`脚本内容如下

```shell
#!/bin/bash
sed -i '/vm.dirty_background_ratio/d' /etc/sysctl.conf
sed -i '/vm.dirty_ratio/d' /etc/sysctl.conf
echo "vm.dirty_ratio=10" >> /etc/sysctl.conf
echo "vm.dirty_background_ratio=5" >> /etc/sysctl.conf
sysctl -p
```

> 5.主机: `10.10.10.1`、`10.10.10.3`、`10.10.10.4`执行`shell`脚本，使用密码登录方式，且`ssh`端口及密码均不一致

```yaml
server:
  - host: 10.10.10.1
    username: root
    # privateKeyPath: ~/.ssh/id_rsa
    password: "123456"
    port: 22
  - host: 10.10.10.3
    username: root
    # privateKeyPath: ~/.ssh/id_rsa
    password: "123"
    port: 22122
  - host: 10.10.10.4
    username: root
    # privateKeyPath: ~/.ssh/id_rsa
    password: "456"
    port: 22222
excludes:
  - 10.10.10.2
script: "/root/modify-sysctl.sh"
```

`/root/modify-sysctl.sh`脚本内容如下

```shell
#!/bin/bash
sed -i '/vm.dirty_background_ratio/d' /etc/sysctl.conf
sed -i '/vm.dirty_ratio/d' /etc/sysctl.conf
echo "vm.dirty_ratio=10" >> /etc/sysctl.conf
echo "vm.dirty_background_ratio=5" >> /etc/sysctl.conf
sysctl -p
```