##  清理dns配置

根据入参操作`/etc/resolv.conf`

### 版本&兼容性

> 版本支持

- [v0.7.8-alpha以上](https://github.com/weiliang-ms/easyctl/releases/)

> 兼容性

- [x] `CentOS6` 
- [x] `CentOS7`

### 使用方式

> 参考以下链接进行安装

- [安装说明文档](../安装文档/README.md)

> 生成默认配置文件

```shell
$ easyctl clean dns
I1001 11:13:06.384839  126576 track.go:50] 检测到配置文件参数为空，生成配置文件样例 -> config.yaml
```

> 修改配置文件

`config.yaml`, 参考[配置样例](#配置样例) 调整配置

```yaml
server:
  - host: 10.10.10.[1:40]
    username: root
    privateKeyPath: "" # ~/.ssh/id_rsa，为空默认走password登录；不为空默认走密钥登录
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
clean-dns:
  address-list:        # 地址列表，为空表示清除所有
    - 8.8.8.8
  excludes:             # 排除哪些dns地址不被清理
    - 114.114.114.114
```

> 执行

```shell
$ easyctl clean dns -c config.yaml --debug
```

### 配置样例

> 1.清空`dns`列表

```yaml
server:
  - host: 10.10.10.[1:40]
    username: root
    privateKeyPath: "" # ~/.ssh/id_rsa，为空默认走password登录；不为空默认走密钥登录
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
clean-dns:
  address-list:        # 地址列表，为空表示清除所有
  excludes:             # 排除哪些dns地址不被清理
```

> 2.删除`8.8.8.8`、`114.114.114.114`

```yaml
server:
  - host: 10.10.10.[1:40]
    username: root
    privateKeyPath: "" # ~/.ssh/id_rsa，为空默认走password登录；不为空默认走密钥登录
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
clean-dns:
  address-list:        # 地址列表，为空表示清除所有
    - 8.8.8.8
    - 114.114.114.114
  excludes:             # 排除哪些dns地址不被清理
```

> 3.清空`dns`配置保留`114.114.114.114`

```yaml
server:
  - host: 10.10.10.[1:40]
    username: root
    privateKeyPath: "" # ~/.ssh/id_rsa，为空默认走password登录；不为空默认走密钥登录
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
clean-dns:
  address-list:        # 地址列表，为空表示清除所有
  excludes:             # 排除哪些dns地址不被清理
  - 114.114.114.114
```